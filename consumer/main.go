package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

const (
	exchangeName       = "watt-flow"
	measurementsBucket = "power_measurements"
	deviceStatusBucket = "device_status"
	influxOrg          = "watt-flow"
	redisTTL           = time.Second * 60
	reconnectDelay     = time.Second * 10
)

type Consumer struct {
	influxClient  influxdb2.Client
	conn          *amqp.Connection
	channel       *amqp.Channel
	pgDB          *sql.DB
	redisClient   *redis.Client
	reconnecting  chan bool
	shutdown      chan struct{}
	cancelContext context.CancelFunc
	amqpURI       string
	wg            sync.WaitGroup
	wsServer      *WsServer
}

type DeviceStatus struct {
	LastSeen time.Time
	DeviceID string
	IsActive bool
}

func NewConsumer(amqpURI, influxURI, influxToken, pgConnStr, redisAddr string, cancelFunc context.CancelFunc, wsServer *WsServer) (*Consumer, error) {
	// Connect to InfluxDB
	influxClient := influxdb2.NewClient(influxURI, influxToken)

	// Connect to PostgreSQL
	pgDB, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		influxClient.Close()
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	err = pgDB.Ping()
	if err != nil {
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
		shutdown:      make(chan struct{}),
		reconnecting:  make(chan bool),
		amqpURI:       amqpURI,
		influxClient:  influxClient,
		pgDB:          pgDB,
		redisClient:   redisClient,
		cancelContext: cancelFunc,
		wsServer:      wsServer,
	}, nil
}

func (c *Consumer) ConnectToBroker() error {
	conn, err := amqp.Dial(c.amqpURI)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %v", err)
	}

	c.conn = conn
	c.channel = channel

	err = c.channel.ExchangeDeclare(
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

	fmt.Println("Successfully initialized rabbitmq connection and exchange!")
	return nil
}

func (c *Consumer) Start(ctx context.Context) error {
	measurementMsgs, err := c.channel.Consume(
		"measurements_queue",
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
		"heartbeats_queue",
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
	c.wg.Add(3)
	go c.processMeasurements(ctx, measurementMsgs)
	go c.processHeartbeats(ctx, heartbeatMsgs)
	go c.updateDeviceStatus(ctx)

	return nil
}

func (c *Consumer) processMeasurements(ctx context.Context, msgs <-chan amqp.Delivery) {
	defer c.wg.Done()
	writeAPI := c.influxClient.WriteAPIBlocking(influxOrg, measurementsBucket)

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
				// log.Printf("Successfully written measurement to InfluxDB")
			}

			// Send consumption WebSocket message for device-specific consumption
			consumptionMsg, err := json.Marshal(Consumption{
				DeviceId:    measurement.DeviceID,
				Consumption: measurement.Value,
			})
			if err != nil {
				log.Printf("Failed marshalling consumption to json!", err)
			} else {
				c.wsServer.SendMessage(measurement.DeviceID, consumptionMsg, "consumption")
				log.Printf("Sent realtime consumption data for device %s: %f kWh", measurement.DeviceID, measurement.Value)
			}

			// Update Redis for city-level consumption aggregation
			cityKey := "city:" + measurement.Address.City
			err = c.redisClient.HIncrByFloat(ctx, cityKey, "value", measurement.Value).Err()
			if err != nil {
				log.Printf("Failed to update Redis: %v", err)
			}
		case <-c.channel.NotifyClose(make(chan *amqp.Error)):
			log.Println("Connection to RabbitMQ lost. Reconnecting...")
			if err := c.reconnectBroker(); err != nil {
				log.Printf("Failed to reconnect to RabbitMQ: %v", err)
			}
		}
	}
}

func debugRedisData(ctx context.Context, redisClient *redis.Client) {
	keys, err := redisClient.Keys(ctx, "city:*").Result()
	if err != nil {
		log.Printf("Failed to fetch Redis keys: %v", err)
		return
	}

	for _, key := range keys {
		value, err := redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to fetch Redis value for key %s: %v", key, err)
			continue
		}

		// Log the key and its value
		log.Printf("Redis Key: %s, Value: %v", key, value)
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
			// err := c.redisClient.Set(ctx, heartbeat.DeviceID, heartbeat.Timestamp, redisTTL).Err()
			err := c.redisClient.HSet(ctx, "heartbeat:"+heartbeat.DeviceID, map[string]interface{}{
				"lastSeen": heartbeat.Timestamp,
			}).Err()
			if err != nil {
				log.Printf("Failed writing to redis!: %v", err)
			}
			c.redisClient.Expire(ctx, "heartbeat:"+heartbeat.DeviceID, redisTTL)

		case <-c.channel.NotifyClose(make(chan *amqp.Error)):
			log.Println("Connection to RabbitMQ lost. Reconnecting...")
			if err := c.reconnectBroker(); err != nil {
				log.Printf("Failed to reconnect to RabbitMQ: %v", err)
			}
		}
	}
}

func (c *Consumer) updateDeviceStatus(ctx context.Context) {
	defer c.wg.Done()
	writeAPI := c.influxClient.WriteAPIBlocking(influxOrg, deviceStatusBucket)

	done := false
	doneTimeout := 0
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.reconnecting:
			log.Println("Cannot process device status! Waiting for connection!")
			time.Sleep(20 * time.Second)
			continue

		case <-ctx.Done():
			log.Printf("Context cancelled. Shutting down!")
			return

		case <-ticker.C:
			keys, err := c.redisClient.Keys(ctx, "heartbeat:*").Result()
			if err != nil {
				log.Printf("Failed to get device IDs from Redis: %v", err)
				continue
			}

			// Iterate through devices and update status
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
				// Checking device status
				is_alive := time.Since(lastSeenTime) <= 30*time.Second

				// If was not previously processed (new device)
				if !exists {
					c.updateStatusInDB(strings.TrimPrefix(key, "heartbeat:"), is_alive, &ctx, writeAPI)
					log.Printf("Updated value in postgres on new entry!")
				} else {
					// only update postgres on status change
					if is_alive != was_alive {
						c.updateStatusInDB(strings.TrimPrefix(key, "heartbeat:"), is_alive, &ctx, writeAPI)
						log.Printf("Updated value in postgres on status change!")
					}
				}
				// update status in redis db
				err = c.redisClient.HSet(ctx, key, "lastStatus", is_alive).Err()
				if err != nil {
					log.Printf("Failed writing to redis!: %v", err)
				}

			}
			if done {
				if doneTimeout >= 8 {
					c.cancelContext()
					return
				}
				log.Printf("Shutdown running! Waiting for possible device status changes.", doneTimeout)
				doneTimeout++
			}

		case <-c.shutdown:
			done = true
		}
	}
}

func (c *Consumer) updateStatusInDB(key string, status bool, ctx *context.Context, writeAPI api.WriteAPIBlocking) {
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

	p := influxdb2.NewPoint(
		"online_status",
		map[string]string{
			"device_id": key,
		},
		map[string]interface{}{
			"value": status,
		},
		time.Now(),
	)

	msg, err := json.Marshal(Status{
		DeviceId: key,
		IsActive: status,
	})
	if err != nil {
		log.Printf("Failed marshalling to json!", err)
	}
	c.wsServer.SendMessage(key, msg, "avb")

	if err := writeAPI.WritePoint(*ctx, p); err != nil {
		log.Printf("Failed to write status change to InfluxDB: %v", err)
	} else {
		log.Printf("Successfully written device online status change to InfluxDB")
	}
}

func (c *Consumer) reconnectBroker() error {
	// Close existing connection and channel
	if err := c.channel.Close(); err != nil {
		log.Printf("Error closing channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
	c.reconnecting <- true

	// Reconnect to RabbitMQ
	for {
		err := c.ConnectToBroker()
		log.Println("Attempting to reconnect to RabbitMQ...")
		if err == nil {
			log.Println("Reconnected to RabbitMQ")
			c.reconnecting <- false
			return nil
		}
		time.Sleep(reconnectDelay)
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
	c.wsServer.Shutdown()
	if err := c.redisClient.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	return nil
}

func (c *Consumer) SendAggregatedMeasurements(ctx context.Context) {
	keys, err := c.redisClient.Keys(ctx, "city:*").Result()
	if err != nil {
		log.Printf("Failed to fetch keys from Redis: %v", err)
		return
	}

	for _, key := range keys {
		data, err := c.redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to fetch data for %s: %v", key, err)
			continue
		}

		city := strings.TrimPrefix(key, "city:")
		value, _ := strconv.ParseFloat(data["value"], 64)

		if value == 0 {
			continue
		}
		debugRedisData(ctx, c.redisClient)
		wsmsg, err := json.Marshal(MeasurementValue{
			Value:     value,
			Timestamp: time.Now(),
			City:      city,
		})
		if err != nil {
			log.Printf("Failed to marshal JSON: %v", err)
			continue
		}
		c.wsServer.SendMessage(city, wsmsg, "csm")

		// reset accumulated value
		c.redisClient.HSet(ctx, key, "value", 0)
	}
}

func main() {
	amqpURI := os.Getenv("AMQP_URI")
	influxURI := os.Getenv("INFLUX_URI")
	influxToken := os.Getenv("INFLUX_TOKEN")
	redisAddr := os.Getenv("REDIS_ADDR")
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_PORT := os.Getenv("DB_PORT")
	DB_TABLE := os.Getenv("DB_TABLE")
	// amqpURI := flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	// influxURI := flag.String("influx", "http://localhost:8086", "InfluxDB URI")
	// influxToken := flag.String("token", "", "InfluxDB token")
	// redisAddr := flag.String("redis", "localhost:6379", "Redis address")
	// flag.Parse()

	pgConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
		DB_HOST,
		DB_PORT,
		DB_USER,
		DB_PASS,
		DB_TABLE,
	)

	ctx, cancel := context.WithCancel(context.Background())

	wsServer := NewWsServer()
	go wsServer.Run()

	go HttpServer(":9000", wsServer)

	consumer, err := NewConsumer(amqpURI, influxURI, influxToken, pgConnStr, redisAddr, cancel, wsServer)
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.ConnectToBroker()
	if err != nil {
		log.Printf("Failed connecting to broker: %v", err)
		os.Exit(1)
	}

	if err := consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}

	//5-min timer for city consumption aggregation
	go func() {
		now := time.Now()
		next := now.Truncate(5 * time.Minute).Add(5 * time.Minute)
		firstWait := time.Until(next)

		firstTimer := time.NewTimer(firstWait)
		defer firstTimer.Stop()
		select {
		case <-ctx.Done():
			return
		case <-firstTimer.C:
			consumer.SendAggregatedMeasurements(ctx)
		}

		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				consumer.SendAggregatedMeasurements(ctx)
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigChan:
			log.Println("Received shutdown signal")
			consumer.shutdown <- struct{}{}
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer shutdownCancel()
			if err := consumer.Shutdown(shutdownCtx); err != nil {
				log.Printf("Error during shutdown: %v", err)
				os.Exit(1)
			}
			return
		}
	}
}
