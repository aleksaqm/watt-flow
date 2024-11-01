package main

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	deviceID  string
	household *Household
)

type Household struct {
	rng            *rand.Rand
	baseLoad       float64
	peakMultiplier float64
	seasonImpact   float64
	address        Location
}

type Location struct {
	countryCode string
	city        string
	street      string
	number      string
}

func newLocation(countryCode string, city string, street string, number string) *Location {
	return &Location{
		countryCode: countryCode,
		city:        city,
		street:      street,
		number:      number,
	}
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

func newHousehold(seed int64, location Location) *Household {
	rng := rand.New(rand.NewSource(seed))
	return &Household{
		rng:            rng,
		baseLoad:       0.3 + (rng.Float64() * 0.2),
		peakMultiplier: 2.5 + rng.Float64(),
		seasonImpact:   0.3 + (rng.Float64() * 0.2),
		address:        location,
	}
}

func (hs *Household) calculateDailyPattern(hour float64) float64 {
	// Morning peak (7-9 AM)
	if hour >= 7 && hour <= 9 {
		return hs.peakMultiplier * 0.8
	}
	// Evening peak (18-22)
	if hour >= 18 && hour <= 22 {
		return hs.peakMultiplier
	}
	// Night time (23-5)
	if hour >= 23 || hour <= 5 {
		return 0.8
	}
	// Mid-day (9-18)
	if hour > 9 && hour < 18 {
		return 1.5
	}
	// Early morning (5-7)
	return 1.2
}

func (hs *Household) calculateSeasonalFactor(month time.Month) float64 {
	monthAngle := float64(month-1) * (2 * math.Pi / 12)
	// Create a sinusoidal pattern with peak in winter (December/January)
	// and trough in summer (June/July)
	seasonalVariation := math.Cos(monthAngle)*hs.seasonImpact + 1.0
	return seasonalVariation
}

func (hs *Household) SimulateConsumption(timestamp time.Time) float64 {
	// Time-of-day factor
	hour := float64(timestamp.Hour())
	dailyPattern := hs.calculateDailyPattern(hour)

	// Seasonal factor
	month := timestamp.Month()
	seasonalFactor := hs.calculateSeasonalFactor(month)

	// Random variation (±15% to simulate appliance usage)
	randomFactor := 0.85 + (hs.rng.Float64() * 0.3)

	// Calculate total consumption
	consumption := hs.baseLoad * dailyPattern * seasonalFactor * randomFactor

	return consumption
}

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

func sendHeartbeat(ctx *context.Context, ch *amqp.Channel) {
	for {
		err := ch.PublishWithContext(*ctx,
			"",
			"heartbeat", // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("heartbeat"),
			})
		failOnError(err, "Failed sending heartbeat!")
		log.Printf(" [x] Sent heartbeat")
		time.Sleep(5 * time.Second)
	}
}

func sendMeasurement(ctx *context.Context, ch *amqp.Channel) {
	startTime := time.Now()
	for {
		value := household.SimulateConsumption(startTime)
		body := newMeasurement(deviceID, value, startTime)
		message, err := json.Marshal(*body)
		failOnError(err, "Failed converting to json")
		fmt.Println(string(message))
		err = ch.PublishWithContext(*ctx,
			"",
			// household.address.city, // routing key
			"measurement",
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(message),
			})
		failOnError(err, "Failed sending measurement!")
		log.Printf(" [x] Sent measurement")
		time.Sleep(1 * time.Minute)
		startTime = startTime.Add(1 * time.Hour)
	}
}

func main() {
	deviceID = *flag.String("device", "be781b42-c3b0-475b-bdc5-cb467d0f7f1b", "uuid of the device")
	country := *flag.String("country", "0", "Country")
	city := *flag.String("city", "0", "city")
	street := *flag.String("street", "0", "street")
	number := *flag.String("number", "0", "street number")

	address := newLocation(country, city, street, number)
	fmt.Println(deviceID)
	deviceIDInt, _ := uuidStrToInt64(deviceID)

	household = newHousehold(deviceIDInt, *address)

	config := amqp.Config{
		Heartbeat: 10 * time.Second,
	}
	conn, err := amqp.DialConfig("amqp://guest:guest@localhost:5672/", config)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = ch.QueueDeclare(
		"heartbeat",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed declaring heartbeat queue")

	_, err = ch.QueueDeclare(
		"measurement",
		false,
		false,
		false,
		false,
		nil)

	failOnError(err, "Failed declaring heartbeat queue")
	go sendHeartbeat(&ctx, ch)
	sendMeasurement(&ctx, ch)
}
