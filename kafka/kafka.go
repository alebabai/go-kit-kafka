package kafka

import (
	"context"
	"time"
)

type Message interface {
	Topic() string
	Partition() int32
	Offset() int64
	Key() []byte
	Value() []byte
	Headers() []Header
	Timestamp() time.Time
}

type Header interface {
	Key() []byte
	Value() []byte
}

type Handler interface {
	Handle(ctx context.Context, msg Message) error
}
