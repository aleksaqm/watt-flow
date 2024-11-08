package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	reconnecting chan struct{}
	exchange     string
	mu           sync.Mutex
	isConnected  bool
}

func NewMessageBroker(exchange string) *MessageBroker {
	return &MessageBroker{
		reconnecting: make(chan struct{}),
		exchange:     exchange,
	}
}

func (b *MessageBroker) Connect() error {
	config := amqp.Config{
		Heartbeat: 10 * time.Second,
	}

	conn, err := amqp.DialConfig("amqp://guest:guest@localhost:5672/", config)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		b.exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	measurementQueue, err := ch.QueueDeclare(
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

	err = ch.QueueBind(
		measurementQueue.Name,
		"measurement.*",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind measurement queue: %v", err)
	}

	// heartbeatQueue, err := ch.QueueDeclare(
	// 	"heartbeats_queue",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to declare heartbeat queue: %v", err)
	// }
	//
	// err = ch.QueueBind(
	// 	heartbeatQueue.Name,
	// 	"heartbeat.*",
	// 	exchangeName,
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to bind heartbeat queue: %v", err)
	// }

	b.mu.Lock()
	b.conn = conn
	b.channel = ch
	b.isConnected = true
	b.mu.Unlock()

	go b.handleConnectionClose()

	return nil
}

func (mb *MessageBroker) handleConnectionClose() {
	<-mb.conn.NotifyClose(make(chan *amqp.Error))
	mb.mu.Lock()
	mb.isConnected = false
	mb.mu.Unlock()

	log.Println("Connection to RabbitMQ lost.")
	for {
		err := mb.Connect()
		log.Println("Attempting to reconnect to RabbitMQ...")
		if err == nil {
			log.Println("Reconnected to RabbitMQ")
			mb.reconnecting <- struct{}{}
			return
		}
		time.Sleep(reconnectDelay)
	}
}

func (mb *MessageBroker) PublishMessage(ctx context.Context, msg Message) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if !mb.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	return mb.channel.PublishWithContext(ctx,
		exchangeName,
		msg.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.Payload,
			Timestamp:   time.Now(),
		})
}
