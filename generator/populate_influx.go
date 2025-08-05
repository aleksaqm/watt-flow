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

	"github.com/influxdata/influxdb-client-go/v2/domain"
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
	timeoutDuration      = 3000 * time.Second
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

	// 2. Fast bucket reset (drop and recreate)
		fmt.Printf("WARNING: Resetting bucket '%s' (drop and recreate)...\n", influxBucket)
		err := resetBucketFast(client, influxOrg, influxBucket)
		if err != nil {
			log.Fatalf("Failed to reset bucket: %v", err)
		}
		fmt.Printf("Successfully reset bucket '%s'.\n", influxBucket)

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

// Fast bucket reset by dropping and recreating
func resetBucketFast(client influxdb2.Client, org, bucket string) error {
	bucketsAPI := client.BucketsAPI()
	ctx := context.Background()

	// 1. Find the existing bucket
	existingBucket, err := bucketsAPI.FindBucketByName(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to find bucket '%s': %w", bucket, err)
	}

	if existingBucket != nil {
		fmt.Printf("Deleting existing bucket '%s'...\n", bucket)
		// 2. Delete the bucket entirely
		err = bucketsAPI.DeleteBucket(ctx, existingBucket)
		if err != nil {
			return fmt.Errorf("failed to delete bucket '%s': %w", bucket, err)
		}
	}

	// 3. Recreate the bucket with same retention policy
	fmt.Printf("Creating new bucket '%s'...\n", bucket)
	orgObj, err := client.OrganizationsAPI().FindOrganizationByName(ctx, org)
	if err != nil {
		return fmt.Errorf("failed to find organization '%s': %w", org, err)
	}

	// Create bucket with default retention (30 days) - adjust as needed
	newBucket := &domain.Bucket{
		Name:           bucket,
		OrgID:          orgObj.Id,
		RetentionRules: []domain.RetentionRule{{EverySeconds: 0}}, // 0 = infinite retention
	}

	_, err = bucketsAPI.CreateBucket(ctx, newBucket)
	if err != nil {
		return fmt.Errorf("failed to create bucket '%s': %w", bucket, err)
	}

	return nil
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
