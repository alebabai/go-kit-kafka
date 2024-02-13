package adapter

import (
	"context"
	"fmt"

	kitkafka "github.com/alebabai/go-kit-kafka/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}

type Producer struct {
	producer producer
}

func NewProducer(producer producer) *Producer {
	return &Producer{
		producer: producer,
	}
}

func (p *Producer) Handle(ctx context.Context, msg *kitkafka.Message) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to produce message: %w", ctx.Err())
	default:
		if err := p.producer.Produce(TransformMessage(msg), nil); err != nil {
			return err
		}

		return nil
	}
}
