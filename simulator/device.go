package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var startTime = time.Now().Add(8 * time.Hour)

type Device struct {
	ID           string
	Household    *Household
	broker       *MessageBroker
	buffer       *MessageBuffer
	logger       *lumberjack.Logger
	lastRotation time.Time
	wg           sync.WaitGroup
}

func NewDevice(deviceID string, household *Household, exchangeName string) (*Device, error) {
	logsDir := "measurements"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Setup lumberjack logger
	logger := &lumberjack.Logger{
		Filename:  filepath.Join(logsDir, "measurements.log"),
		MaxAge:    90, // days
		Compress:  false,
		LocalTime: false,
	}
	return &Device{
		ID:           deviceID,
		Household:    household,
		broker:       NewMessageBroker(exchangeName),
		buffer:       NewMessageBuffer(),
		logger:       logger,
		lastRotation: startTime,
	}, nil
}

func (d *Device) Start(ctx context.Context) error {
	if err := d.broker.Connect(); err != nil {
		return err
	}

	d.wg.Add(2)
	go func() {
		defer d.wg.Done()
		d.sendHeartbeats(ctx)
	}()
	go func() {
		defer d.wg.Done()
		d.sendMeasurements(ctx, startTime)
	}()
	return nil
}
func (d *Device) Shutdown(ctx context.Context) error {
	if err := d.logger.Close(); err != nil {
		return fmt.Errorf("failed to close logger: %v", err)
	}

	if d.broker.conn != nil {
		if err := d.broker.conn.Close(); err != nil {
			return fmt.Errorf("failed to close broker connection: %v", err)
		}
	}

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timed out: %v", ctx.Err())
	}
}

func (d *Device) sendHeartbeats(ctx context.Context) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			heartbeat := newHeartbeat(d.ID, time.Now())
			payload, _ := json.Marshal(heartbeat)
			msg := Message{
				Type:      "heartbeat",
				Payload:   payload,
				Queue:     "heartbeat." + d.Household.address.city,
				Timestamp: time.Now(),
			}

			if err := d.broker.PublishMessage(ctx, msg); err != nil {
				log.Printf("Failed to send heartbeat: %v", err)
			} else {
				log.Printf("Sent heartbeat")
			}
		}
	}
}

func (d *Device) sendMeasurements(ctx context.Context, currentTime time.Time) {
	ticker := time.NewTicker(measurementInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-d.broker.reconnecting:
			// On reconnect, try to send buffered messages
			d.sendBufferedMessages(ctx)
		case <-ticker.C:

			if err := d.checkAndRotateLog(currentTime); err != nil {
				log.Printf("Error rotating log: %v", err)
			}

			value := d.Household.SimulateConsumption(currentTime)
			measurement := newMeasurement(d.ID, value, currentTime)
			payload, _ := json.Marshal(measurement)

			msg := Message{
				Type:      "measurement",
				Payload:   payload,
				Queue:     "measurement." + d.Household.address.city,
				Timestamp: currentTime,
			}

			if err := d.logMeasurement(measurement); err != nil {
				log.Printf("Failed to log measurement: %v", err)
			}
			// Try to send measurement
			if err := d.broker.PublishMessage(ctx, msg); err != nil {
				log.Printf("Failed to send measurement: %v", err)
				d.buffer.Add(msg)
			} else {
				log.Printf("Sent measurement")
			}

			currentTime = currentTime.Add(1 * time.Hour)
		}
	}
}

func (d *Device) sendBufferedMessages(ctx context.Context) {
	messages := d.buffer.Flush()
	for _, msg := range messages {
		if err := d.broker.PublishMessage(ctx, msg); err != nil {
			// If still failing, add back to buffer
			d.buffer.Add(msg)
		} else {
			log.Printf("Sent buffered message of type: %s", msg.Type)
		}
	}
}
func (d *Device) checkAndRotateLog(currentTime time.Time) error {
	if currentTime.Day() != d.lastRotation.Day() {
		//d.logger.Filename = filepath.Join("logs", fmt.Sprintf("measurements_%s.log", currentTime.Format("2006-01-02")))

		// Rotate the log file
		if err := d.logger.Rotate(); err != nil {
			return fmt.Errorf("failed to rotate log file: %v", err)
		}

		d.lastRotation = currentTime
		log.Printf("Rotated log file to %s", d.logger.Filename)
	}
	return nil
}

func (d *Device) logMeasurement(m *Measurement) error {
	logEntry := struct {
		Timestamp time.Time `json:"timestamp"`
		DeviceID  string    `json:"device_id"`
		Value     float64   `json:"value"`
	}{
		Timestamp: m.Timestamp,
		DeviceID:  m.DeviceID,
		Value:     m.Value,
	}

	data, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	_, err = d.logger.Write(append(data, '\n'))
	return err
}
