package main

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	reconnectDelay       = 5 * time.Second
	heartbeatInterval    = 5 * time.Second
	measurementInterval  = 1 * time.Minute
	measurementRetention = 90 * 24 * time.Hour
	exchangeName         = "watt-flow"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func uuidStrToInt64(uStr string) (int64, error) {
	u, err := uuid.Parse(strings.TrimSpace(uStr))
	if err != nil {
		log.Panic("Invalid device ID")
	}
	hash := md5.Sum(u[:])
	return int64(binary.BigEndian.Uint64(hash[:8])), nil
}

func main() {
	deviceID := flag.String("device", "0", "uuid of the device")
	country := flag.String("country", "0", "Country")
	city := flag.String("city", "0", "city")
	street := flag.String("street", "0", "street")
	number := flag.String("number", "0", "street number")
	flag.Parse()

	if *city == "0" || *country == "0" {
		log.Fatal("city and country args are required!")
	}

	address := newLocation(*country, *city, *street, *number)
	deviceIDInt, err := uuidStrToInt64(*deviceID)
	if err != nil {
		log.Fatal(err)
	}

	household := newHousehold(deviceIDInt, *address)
	device, err := NewDevice(*deviceID, household, exchangeName)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := device.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// Keep the main goroutine running
	select {}
}
