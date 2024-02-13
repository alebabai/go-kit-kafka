package kafka

import (
	"context"
	"time"
)

// Handler is an interface for processing Kafka messages.
type Handler interface {
	Handle(ctx context.Context, msg *Message) error
}

// Message is a Kafka message object.
type Message struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
	Headers   []Header
	Timestamp time.Time
}

// Header is a Kafka message object.
type Header struct {
	Key   []byte
	Value []byte
}
