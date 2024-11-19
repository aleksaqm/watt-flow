package main

import "time"

type Heartbeat struct {
	DeviceID  string
	Timestamp string
}

type Measurement struct {
	DeviceID  string
	Value     float64
	Timestamp time.Time
	Address   Location
}

type Location struct {
	City   string
	Street string
	Number string
}

type Status struct {
	DeviceId string
	IsActive bool
}
