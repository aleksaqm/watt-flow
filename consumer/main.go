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
	"strconv"
	"sync"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

const (
	exchangeName = "watt-flow"
	influxBucket = "power_measurements"
	influxOrg    = "watt-flow"
	redisTTL     = time.Second * 60
)

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	influxClient influxdb2.Client
	pgDB         *sql.DB
	redisClient  *redis.Client
	wg           sync.WaitGroup
}

type DeviceStatus struct {
	LastSeen time.Time
	DeviceID string
	IsActive bool
}

func NewConsumer(amqpURI, influxURI, influxToken, pgConnStr, redisAddr string) (*Consumer, error) {
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
	influxClient := influxdb2.NewClient(influxURI, influxToken)

	// Connect to PostgreSQL
	pgDB, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		channel.Close()
		conn.Close()
		influxClient.Close()
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	err = pgDB.Ping()
	if err != nil {
		channel.Close()
		conn.Close()
		influxClient.Close()
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}
	fmt.Println("Successfully connected to postgres!")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "password", // no password set
		DB:       0,          // use default DB
	})

	// Initialize PostgreSQL table
	_, err = pgDB.Exec(`
        CREATE TABLE IF NOT EXISTS device_status (
            device_id TEXT PRIMARY KEY,
            is_active BOOLEAN
        )
    `)
	if err != nil {
		fmt.Printf("Failed creating table: %v", err)
	}

	return &Consumer{
		conn:         conn,
		channel:      channel,
		influxClient: influxClient,
		pgDB:         pgDB,
		redisClient:  redisClient,
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
		false,
		true,
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

	fmt.Println("Successfully initialized rabbitmq connection and exchange!")

	c.wg.Add(3)
	go c.processMeasurements(ctx, measurementMsgs)
	go c.processHeartbeats(ctx, heartbeatMsgs)
	go c.updateDeviceStatus(ctx)

	return nil
}

func (c *Consumer) processMeasurements(ctx context.Context, msgs <-chan amqp.Delivery) {
	defer c.wg.Done()
	writeAPI := c.influxClient.WriteAPIBlocking(influxOrg, influxBucket)

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

			p := influxdb2.NewPoint(
				"power_consumption",
				map[string]string{
					"device_id": measurement.DeviceID,
					"city":      measurement.Address.City,
				},
				map[string]interface{}{
					"value": measurement.Value,
				},
				measurement.Timestamp,
			)

			if err := writeAPI.WritePoint(ctx, p); err != nil {
				log.Printf("Failed to write measurement to InfluxDB: %v", err)
			} else {
				log.Printf("Successfully written measurement to InfluxDB")
			}
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
			// err := c.redisClient.Set(ctx, heartbeat.DeviceID, heartbeat.Timestamp, redisTTL).Err()
			err := c.redisClient.HSet(ctx, heartbeat.DeviceID, map[string]interface{}{
				"lastSeen": heartbeat.Timestamp,
			}).Err()
			if err != nil {
				log.Printf("Failed writing to redis!: %v", err)
			}
			c.redisClient.Expire(ctx, heartbeat.DeviceID, redisTTL)

		}
	}
}

func (c *Consumer) updateDeviceStatus(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			keys, err := c.redisClient.Keys(ctx, "*").Result()
			if err != nil {
				log.Printf("Failed to get device IDs from Redis: %v", err)
				continue
			}

			// Iterate through devices and update status in PostgreSQL
			for _, key := range keys {
				result, err := c.redisClient.HGetAll(ctx, key).Result()
				if err != nil {
					log.Printf("Failed to get data for device %s: %v", key, err)
					continue
				}
				lastStatus, exists := result["lastStatus"]
				lastSeen := result["lastSeen"]
				was_alive, _ := strconv.ParseBool(lastStatus)

				lastSeenTime, err := time.Parse(time.RFC3339, lastSeen)
				if err != nil {
					log.Printf("Failed to parse last seen time for device %s: %v", key, err)
					continue
				}
				is_alive := time.Since(lastSeenTime) <= 30*time.Second

				if !exists {
					c.updateStatusInDB(key, is_alive, &ctx)
					log.Printf("Updated value in posgres on new entry!")
				} else {
					if is_alive != was_alive {
						c.updateStatusInDB(key, is_alive, &ctx)
						log.Printf("Updated value in posgres on status change!")
					}
				}
				err = c.redisClient.HSet(ctx, key, "lastStatus", is_alive).Err()
				if err != nil {
					log.Printf("Failed writing to redis!: %v", err)
				}

			}
		}
	}
}

func (c *Consumer) updateStatusInDB(key string, status bool, ctx *context.Context) {
	_, err := c.pgDB.ExecContext(*ctx,
		`INSERT INTO device_status (device_id, is_active)
                         VALUES ($1, $2)
                         ON CONFLICT (device_id)
                         DO UPDATE SET is_active = $2`,
		key,
		status,
	)
	if err != nil {
		log.Printf("Failed to update device status for %s: %v", key, err)
	}
}

func (c *Consumer) Shutdown(ctx context.Context) error {
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
	if err := c.pgDB.Close(); err != nil {
		log.Printf("Error closing PostgreSQL connection: %v", err)
	}
	c.influxClient.Close()
	if err := c.redisClient.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	return nil
}

func main() {
	amqpURI := flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	influxURI := flag.String("influx", "http://localhost:8086", "InfluxDB URI")
	influxToken := flag.String("token", "", "InfluxDB token")
	// pgConnStr := flag.String("pg", "postgres://postgres:postgres@localhost:5432/watt-flow?sslmode=disable", "PostgreSQL connection string")
	redisAddr := flag.String("redis", "localhost:6379", "Redis address")
	flag.Parse()

	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "watt-flow")
	consumer, err := NewConsumer(*amqpURI, *influxURI, *influxToken, pgConnStr, *redisAddr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err := consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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
		}
	}
}
