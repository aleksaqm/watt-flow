package main

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type Device struct {
	ID        string
	Household *Household
	broker    *MessageBroker
	buffer    *MessageBuffer
}

func NewDevice(deviceID string, household *Household, exchangeName string) (*Device, error) {
	//// Create data directory if it doesn't exist
	//dataDir := "data"
	//if err := os.MkdirAll(dataDir, 0755); err != nil {
	//	return nil, fmt.Errorf("failed to create data directory: %v", err)
	//}
	//
	//// Open LevelDB database
	//db, err := leveldb.OpenFile(filepath.Join(dataDir, "measurements.db"), nil)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to open database: %v", err)
	//}

	return &Device{
		ID:        deviceID,
		Household: household,
		broker:    NewMessageBroker(exchangeName),
		buffer:    NewMessageBuffer(),
	}, nil
}

func (d *Device) Start(ctx context.Context) error {
	if err := d.broker.Connect(); err != nil {
		return err
	}

	// Start heartbeat and measurement routines
	go d.sendHeartbeats(ctx)
	go d.sendMeasurements(ctx)
	//go d.cleanupOldMeasurements(ctx)

	return nil
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

func (d *Device) sendMeasurements(ctx context.Context) {
	ticker := time.NewTicker(measurementInterval)
	defer ticker.Stop()

	startTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case <-d.broker.reconnecting:
			// On reconnect, try to send buffered messages
			d.sendBufferedMessages(ctx)
		case <-ticker.C:
			value := d.Household.SimulateConsumption(startTime)
			measurement := newMeasurement(d.ID, value, startTime)
			payload, _ := json.Marshal(measurement)

			msg := Message{
				Type:      "measurement",
				Payload:   payload,
				Queue:     "measurement." + d.Household.address.city,
				Timestamp: startTime,
			}

			// Store measurement in local database TODO
			//if err := d.storeMeasurement(measurement); err != nil {
			//	log.Printf("Failed to store measurement: %v", err)
			//}

			// Try to send measurement
			if err := d.broker.PublishMessage(ctx, msg); err != nil {
				log.Printf("Failed to send measurement: %v", err)
				d.buffer.Add(msg)
			} else {
				log.Printf("Sent measurement")
			}

			startTime = startTime.Add(1 * time.Hour)
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

//func (d *Device) storeMeasurement(m *Measurement) error {
//	key := fmt.Sprintf("%d", m.Timestamp.Unix())
//	value, err := json.Marshal(m)
//	if err != nil {
//		return err
//	}
//	return d.db.Put([]byte(key), value, nil)
//}
