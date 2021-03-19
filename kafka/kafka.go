package kafka

import (
	"context"
	"io"
	"time"
)

type Message interface {
	Topic() string
	Partition() int
	Offset() int64
	Key() []byte
	Value() []byte
	Headers() []Header
	Timestamp() time.Time
}

type Header interface {
	Key() string
	Value() []byte
}

type Reader interface {
	ReadMessage(ctx context.Context, timeout time.Duration) (Message, error)
	io.Closer
}

type Handler interface {
	Handle(ctx context.Context, msg Message) error
}
