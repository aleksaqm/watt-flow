package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type DeviceOptions struct {
	startDate time.Time
	downtime  time.Duration
}

type Device struct {
	lastRotation time.Time
	Household    *Household
	broker       *MessageBroker
	buffer       *MessageBuffer
	logger       *lumberjack.Logger
	config       *Config
	ID           string
	wg           sync.WaitGroup
}

func NewDevice(deviceID string, household *Household, exchangeName string, amqpUri string) (*Device, error) {
	logsDir := "measurements/" + deviceID
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
	config := &Config{
		filename: "config/config.json",
		data:     &ConfigEntry{},
	}
	return &Device{
		ID:        deviceID,
		Household: household,
		broker:    NewMessageBroker(exchangeName, amqpUri),
		buffer:    NewMessageBuffer(),
		logger:    logger,
		config:    config,
	}, nil
}

func (d *Device) Start(ctx context.Context, options *DeviceOptions) error {
	if err := d.broker.Connect(); err != nil {
		return err
	}
	if options == nil {
		if err := d.config.LoadConfig(); err != nil {
			return err
		}
	} else {
		d.config.data.LastMeasurement = options.startDate
		d.config.data.DowntimeSimulation = options.downtime
	}
	d.lastRotation = d.config.data.LastMeasurement
	d.config.data.LastMeasurement = d.config.data.LastMeasurement.Add(d.config.data.DowntimeSimulation)

	d.wg.Add(2)
	go func() {
		defer d.wg.Done()
		d.sendHeartbeats(ctx)
	}()
	go func() {
		defer d.wg.Done()
		d.sendMeasurements(ctx, d.config.data.LastMeasurement)
	}()
	return nil
}

func (d *Device) Shutdown(ctx context.Context) error {
	if err := d.logger.Rotate(); err != nil {
		return fmt.Errorf("failed to rotate logger: %v", err)
	}
	if err := d.logger.Close(); err != nil {
		return fmt.Errorf("failed to close logger: %v", err)
	}

	if d.broker.conn != nil {
		if err := d.broker.conn.Close(); err != nil {
			return fmt.Errorf("failed to close broker connection: %v", err)
		}
	}

	if err := d.config.SaveConfig(); err != nil {
		log.Printf("failed to save config: %v", err)
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
			heartbeat := newHeartbeat(d.ID, time.Now().Format(time.RFC3339))
			payload, _ := json.Marshal(heartbeat)
			msg := Message{
				Type:      "heartbeat",
				Payload:   payload,
				Queue:     "heartbeat." + d.Household.address.City,
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
			d.config.data.LastMeasurement = currentTime
			if err := d.checkAndRotateLog(currentTime); err != nil {
				log.Printf("Error rotating log: %v", err)
			}

			value := d.Household.SimulateConsumption(currentTime)
			measurement := newMeasurement(d.ID, value, currentTime, *d.Household.address)
			payload, _ := json.Marshal(measurement)
			msg := Message{
				Type:      "measurement",
				Payload:   payload,
				Queue:     "measurement." + d.Household.address.City,
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
		// d.logger.Filename = filepath.Join("logs", fmt.Sprintf("measurements_%s.log", currentTime.Format("2006-01-02")))

		// Rotate the log file
		if err := d.logger.Rotate(); err != nil {
			return fmt.Errorf("failed to rotate log file: %v", err)
		}
		if err := d.config.SaveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		d.lastRotation = currentTime
		log.Printf("Rotated log file to %s", d.logger.Filename)
	}
	return nil
}

func (d *Device) logMeasurement(m *Measurement) error {
	logEntry := m.Timestamp.Local().String() + "->" + strconv.FormatFloat(m.Value, 'f', -1, 64)

	_, err := d.logger.Write(append([]byte(logEntry), '\n'))
	return err
}
