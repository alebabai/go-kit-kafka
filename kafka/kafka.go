package kafka

import (
	"context"
	"time"
)

// Message represents a Kafka message
type Message struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
	Headers   []Header
	Timestamp time.Time
}

// Header represents a Kafka header
type Header struct {
	Key   []byte
	Value []byte
}

// Handler represents a Kafka message handler
type Handler interface {
	Handle(ctx context.Context, msg *Message) error
}
