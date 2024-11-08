package main

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
)

const (
	reconnectDelay      = 5 * time.Second
	heartbeatInterval   = 5 * time.Second
	measurementInterval = 1 * time.Minute
	exchangeName        = "watt-flow"
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
	deviceID := os.Getenv("DEVICE_ID")
	city := os.Getenv("DEVICE_CITY")
	street := os.Getenv("DEVICE_STREET")
	number := os.Getenv("DEVICE_NUMBER")
	amqpURI := os.Getenv("AMQP_URI")
	// deviceID := flag.String("device", "0", "uuid of the device")
	// city := flag.String("city", "0", "city")
	// street := flag.String("street", "0", "street")
	// number := flag.String("number", "0", "street number")
	// flag.Parse()

	if city == "" || street == "" || number == "" || amqpURI == "" {
		log.Fatal("city and address args are required!")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create context that will be canceled on signal
	ctx, cancel := context.WithCancel(context.Background())

	address := newLocation(city, street, number)
	deviceIDInt, err := uuidStrToInt64(deviceID)
	if err != nil {
		log.Fatal(err)
	}

	household := newHousehold(deviceIDInt, address)
	device, err := NewDevice(deviceID, household, exchangeName, amqpURI)
	if err != nil {
		log.Fatal(err)
	}

	if err := device.Start(ctx, nil); err != nil {
		cancel()
		log.Fatal(err)
	}

	// Wait for interrupt signal
	sig := <-sigChan
	log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

	// Cancel context to stop all routines
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Perform graceful shutdown
	if err := device.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Graceful shutdown completed")
}
