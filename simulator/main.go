package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
		"measurment",
		false,
		false,
		false,
		false,
		nil)

	failOnError(err, "Failed declaring heartbeat queue")

	body := "danilo"
	err = ch.PublishWithContext(ctx,
		"",
		"heartbeat", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
