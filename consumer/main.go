package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "watt-flow"
	influxBucket = "power_measurements"
	influxOrg    = "nvt12"
)

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	influxClient influxdb2.Client
	pgDB         *sql.DB
	wg           sync.WaitGroup
}

type DeviceStatus struct {
	LastSeen time.Time
	DeviceID string
	IsActive bool
}

func NewConsumer(amqpURI, influxURI, influxToken, pgConnStr string) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Connect to InfluxDB
	//influxClient := influxdb2.NewClient(influxURI, influxToken)
	//
	//// Connect to PostgreSQL
	//pgDB, err := sql.Open("postgres", pgConnStr)
	//if err != nil {
	//	channel.Close()
	//	conn.Close()
	//	influxClient.Close()
	//	return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	//}

	return &Consumer{
		conn:    conn,
		channel: channel,
		// influxClient: influxClient,
		// pgDB:         pgDB,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	err := c.channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Setup measurement queue
	measurementQueue, err := c.channel.QueueDeclare(
		"measurements_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare measurement queue: %v", err)
	}

	err = c.channel.QueueBind(
		measurementQueue.Name,
		"measurement.*",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind measurement queue: %v", err)
	}

	// Setup heartbeat queue
	heartbeatQueue, err := c.channel.QueueDeclare(
		"heartbeats_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare heartbeat queue: %v", err)
	}

	err = c.channel.QueueBind(
		heartbeatQueue.Name,
		"heartbeat.*",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind heartbeat queue: %v", err)
	}

	measurementMsgs, err := c.channel.Consume(
		measurementQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start measurement consumer: %v", err)
	}

	heartbeatMsgs, err := c.channel.Consume(
		heartbeatQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start heartbeat consumer: %v", err)
	}

	c.wg.Add(2)
	go c.processMeasurements(ctx, measurementMsgs)
	go c.processHeartbeats(ctx, heartbeatMsgs)

	return nil
}

func (c *Consumer) processMeasurements(ctx context.Context, msgs <-chan amqp.Delivery) {
	defer c.wg.Done()
	// writeAPI := c.influxClient.WriteAPIBlocking(influxOrg, influxBucket)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			var measurement Measurement
			if err := json.Unmarshal(msg.Body, &measurement); err != nil {
				log.Printf("Failed to unmarshal measurement: %v", err)
				continue
			}
			log.Println(measurement)

			// Write to InfluxDB
			//p := influxdb2.NewPoint(
			//	"power_consumption",
			//	map[string]string{
			//		"device_id": measurement.DeviceID,
			//	},
			//	map[string]interface{}{
			//		"value": measurement.Value,
			//	},
			//	measurement.Timestamp,
			//)
			//
			//if err := writeAPI.WritePoint(ctx, p); err != nil {
			//	log.Printf("Failed to write measurement to InfluxDB: %v", err)
			//}
		}
	}
}

func (c *Consumer) processHeartbeats(ctx context.Context, msgs <-chan amqp.Delivery) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			var heartbeat Heartbeat
			if err := json.Unmarshal(msg.Body, &heartbeat); err != nil {
				log.Printf("Failed to unmarshal heartbeat: %v", err)
				continue
			}
			log.Println(heartbeat)

			// Update PostgreSQL
			//_, err := c.pgDB.ExecContext(ctx,
			//	`INSERT INTO device_status (device_id, last_seen, is_active)
			//     VALUES ($1, $2, true)
			//     ON CONFLICT (device_id)
			//     DO UPDATE SET last_seen = $2, is_active = true`,
			//	heartbeat.DeviceID,
			//	heartbeat.Timestamp,
			//)
			//if err != nil {
			//	log.Printf("Failed to update device status: %v", err)
			//}
		}
	}
}

func (c *Consumer) markInactiveDevices(ctx context.Context) error {
	_, err := c.pgDB.ExecContext(ctx,
		`UPDATE device_status 
         SET is_active = false 
         WHERE last_seen < $1`,
		time.Now().Add(-15*time.Second),
	)
	return err
}

func (c *Consumer) Shutdown(ctx context.Context) error {
	// Cancel context will stop the processing goroutines
	// Wait for goroutines to finish
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return fmt.Errorf("shutdown timed out: %v", ctx.Err())
	}

	// Close connections
	if err := c.channel.Close(); err != nil {
		log.Printf("Error closing channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
	//if err := c.pgDB.Close(); err != nil {
	//	log.Printf("Error closing PostgreSQL connection: %v", err)
	//}
	//c.influxClient.Close()

	return nil
}

func main() {
	amqpURI := flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	influxURI := flag.String("influx", "http://localhost:8086", "InfluxDB URI")
	influxToken := flag.String("token", "", "InfluxDB token")
	pgConnStr := flag.String("pg", "postgres://username:password@localhost:5432/dbname?sslmode=disable", "PostgreSQL connection string")
	flag.Parse()

	consumer, err := NewConsumer(*amqpURI, *influxURI, *influxToken, *pgConnStr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err := consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			log.Println("Received shutdown signal")
			cancel()
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()
			if err := consumer.Shutdown(shutdownCtx); err != nil {
				log.Printf("Error during shutdown: %v", err)
				os.Exit(1)
			}
			return
		case <-ticker.C:
			//if err := consumer.markInactiveDevices(ctx); err != nil {
			//	log.Printf("Error marking inactive devices: %v", err)
			//}
		}
	}
}
