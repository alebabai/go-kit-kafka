package adapter

import (
	"context"
	"io"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	kitkafka "github.com/alebabai/go-kit-kafka/kafka"
)

type consumer interface {
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Commit() ([]kafka.TopicPartition, error)
	io.Closer
}

type Reader struct {
	consumer consumer
}

func NewReader(c consumer) (*Reader, error) {
	r := &Reader{
		consumer: c,
	}

	return r, nil
}

func (r *Reader) ReadMessage(_ context.Context) (kitkafka.Message, error) {
	msg, err := r.consumer.ReadMessage(-1)
	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

func (r *Reader) Commit(_ context.Context) error {
	if _, err := r.consumer.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Reader) Close() error {
	if r.consumer != nil {
		return r.consumer.Close()
	}

	return nil
}
