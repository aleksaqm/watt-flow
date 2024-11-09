package main

import "sync"

type MessageBuffer struct {
	messages []Message
	mu       sync.Mutex
}

func NewMessageBuffer() *MessageBuffer {
	return &MessageBuffer{
		messages: make([]Message, 0),
	}
}

func (mb *MessageBuffer) Add(msg Message) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.messages = append(mb.messages, msg)
}

func (mb *MessageBuffer) Flush() []Message {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	messages := mb.messages
	mb.messages = make([]Message, 0)
	return messages
}
