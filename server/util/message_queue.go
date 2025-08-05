package util

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewMessageQueue initializes a new RabbitMQ connection
func NewMessageQueue(queueName string) (*MessageQueue, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel to RabbitMQ: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto-delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create a queue on RabbitMQ: %v", err)
	}

	return &MessageQueue{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (mq *MessageQueue) Publish(queueName string, message []byte) error {
	err := mq.channel.Publish(
		"",        // Exchange (empty means default)
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
	return err
}

func (mq *MessageQueue) Close() {
	mq.channel.Close()
	mq.conn.Close()
}
