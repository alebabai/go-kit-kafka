package adapter

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	kitkafka "github.com/alebabai/go-kit-kafka/kafka"
)

type consumer interface {
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Events() chan kafka.Event
	Assign(partitions []kafka.TopicPartition) (err error)
	Unassign() (err error)
	io.Closer
}

type ChannelReader struct {
	consumer consumer
}

func NewChannelReader(c consumer) (*ChannelReader, error) {
	r := &ChannelReader{
		consumer: c,
	}

	return r, nil
}

func (r *ChannelReader) ReadMessage(ctx context.Context, timeout time.Duration) (kitkafka.Message, error) {
	switch e := (<-r.consumer.Events()).(type) {
	case kafka.AssignedPartitions:
		if err := r.consumer.Assign(e.Partitions); err != nil {
			return nil, err
		}
		return nil, errors.New(e.String())
	case kafka.RevokedPartitions:
		if err := r.consumer.Unassign(); err != nil {
			return nil, err
		}
		return nil, errors.New(e.String())
	case *kafka.Message:
		return NewMessage(e), nil
	case kafka.Error:
		return nil, e
	default:
		return nil, errors.New(e.String())
	}
}

func (r *ChannelReader) Close() error {
	if r.consumer != nil {
		return r.consumer.Close()
	}

	return nil
}

type FunctionReader struct {
	consumer consumer
}

func NewFunctionReader(c consumer) (*FunctionReader, error) {
	r := &FunctionReader{
		consumer: c,
	}

	return r, nil
}

func (r *FunctionReader) ReadMessage(ctx context.Context, timeout time.Duration) (kitkafka.Message, error) {
	msg, err := r.consumer.ReadMessage(timeout)
	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

func (r *FunctionReader) Close() error {
	if r.consumer != nil {
		return r.consumer.Close()
	}

	return nil
}
