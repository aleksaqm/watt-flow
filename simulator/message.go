package main

import "time"

type Heartbeat struct {
	DeviceID  string
	Timestamp time.Time
}

func newHeartbeat(device string, timestamp time.Time) *Heartbeat {
	return &Heartbeat{
		DeviceID:  device,
		Timestamp: timestamp,
	}
}

type Message struct {
	Type      string
	Payload   []byte
	Queue     string
	Timestamp time.Time
}

type Measurement struct {
	DeviceID  string
	Value     float64
	Timestamp time.Time
}

func newMeasurement(deviceID string, value float64, timestamp time.Time) *Measurement {
	return &Measurement{
		DeviceID:  deviceID,
		Value:     value,
		Timestamp: timestamp,
	}
}
