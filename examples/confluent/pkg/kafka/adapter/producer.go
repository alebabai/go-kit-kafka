package adapter

import (
	"context"
	"fmt"

	"github.com/alebabai/go-kafka"
	adapter "github.com/alebabai/go-kafka/adapter/confluent"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type producer interface {
	Produce(msg *ckafka.Message, deliveryChan chan ckafka.Event) error
}

type Producer struct {
	producer producer
}

func NewProducer(p producer) *Producer {
	return &Producer{
		producer: p,
	}
}

func (p *Producer) Handle(ctx context.Context, msg kafka.Message) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to produce message: %w", ctx.Err())
	default:
		pmsg := adapter.ConvertKafkaMessageToMessage(msg)
		if err := p.producer.Produce(&pmsg, nil); err != nil {
			return fmt.Errorf("failed to produce message: %w", err)
		}

		return nil
	}
}
