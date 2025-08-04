package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	sim "generator/simulator"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// --- CONFIGURATION ---
// (Copied from your Consumer code)
const (
	influxURI            = "http://localhost:8086"
	influxToken          = "RjRV5kMaSqdQnQ_GlKWqT74BNT_E3R8T0qaiOYiVuLTygokZHZHlaOEmerKGIjS2Gb0Vkbr6v-jR4VtI4mN9qg=="
	influxOrg            = "watt-flow"
	influxBucket         = "power_measurements"
	csvFile              = "simulators.csv"
	historyDuration      = 3 * 365 * 24 * time.Hour
	influxWriteBatchSize = 5000
	timeoutDuration      = 300 * time.Second
)

// Struct to hold data from the CSV file
type Device struct {
	ID   string
	City string
}

func main() {
	fmt.Println("--- InfluxDB Historical Data Backfill Script ---")

	// 1. Create InfluxDB client
	httpClient := &http.Client{
		Timeout: timeoutDuration,
	}
	options := influxdb2.DefaultOptions().
		SetHTTPClient(httpClient)
	client := influxdb2.NewClientWithOptions(influxURI, influxToken, options)
	fmt.Printf("Successfully connected to InfluxDB at %s\n", influxURI)

	// 2. Delete all existing data in the bucket
	fmt.Printf("WARNING: Deleting all data from bucket '%s'...\n", influxBucket)
	err := deleteBucketData(client, influxOrg, influxBucket)
	if err != nil {
		log.Fatalf("Failed to delete bucket data: %v", err)
	}
	fmt.Printf("Successfully cleared bucket '%s'.\n", influxBucket)

	// 3. Read devices from CSV
	devices, err := readDevicesFromCSV(csvFile)
	if err != nil {
		log.Fatalf("Failed to read devices from CSV: %v", err)
	}
	fmt.Printf("Found %d devices in %s to process.\n\n", len(devices), csvFile)

	// 4. Process each device and send data
	// Use the non-blocking WriteAPI for performance
	writeAPI := client.WriteAPI(influxOrg, influxBucket)

	// Process each device sequentially
	for i, device := range devices {
		fmt.Printf("[%d/%d] Processing device: %s (%s)...\n", i+1, len(devices), device.ID, device.City)
		processDevice(writeAPI, device)
	}

	// 5. Flush all buffered writes to InfluxDB
	fmt.Println("\nFlushing all remaining data points to InfluxDB...")
	writeAPI.Flush()
	fmt.Println("--- Backfill complete! ---")
}

// Deletes all data within a specified time range in a bucket.
func deleteBucketData(client influxdb2.Client, org, bucket string) error {
	deleteAPI := client.DeleteAPI()
	// Define a time range that covers all possible data
	start := time.Unix(0, 0)
	stop := time.Now()

	return deleteAPI.DeleteWithName(context.Background(), org, bucket, start, stop, "")
}

// Reads the device_id and city from the specified CSV file.
func readDevicesFromCSV(filePath string) ([]Device, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open csv file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read and discard the header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("could not read header: %w", err)
	}

	var devices []Device
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %w", err)
		}
		if len(record) < 2 {
			log.Printf("Skipping invalid record: %v", record)
			continue
		}
		devices = append(devices, Device{ID: record[0], City: record[1]})
	}
	return devices, nil
}

// Generates and writes 3 years of historical data for a single device.
func processDevice(writeAPI api.WriteAPI, device Device) {
	// Create a deterministic seed from the device ID for reproducible simulations
	hasher := fnv.New64a()
	hasher.Write([]byte(device.ID))
	seed := int64(hasher.Sum64())

	location := sim.NewLocation(device.City, "Simulated Street", "123")
	household := sim.NewHousehold(seed, location)

	endTime := time.Now()
	startTime := endTime.Add(-historyDuration)
	pointsGenerated := 0

	// Loop from 3 years ago until now, in 1-hour increments
	for t := startTime; t.Before(endTime); t = t.Add(time.Hour) {
		consumption := household.SimulateConsumption(t)

		// Create an InfluxDB data point
		p := influxdb2.NewPoint(
			"power_consumption",
			map[string]string{ // Tags
				"device_id": device.ID,
				"city":      device.City,
			},
			map[string]interface{}{ // Fields
				"value": consumption,
			},
			t, // Timestamp
		)

		// Add the point to the batch writer. It will be sent automatically
		// when the batch size is reached or when Flush() is called.
		writeAPI.WritePoint(p)
		pointsGenerated++
	}

	fmt.Printf("  -> Generated and queued %d data points for device %s.\n", pointsGenerated, device.ID)
}
