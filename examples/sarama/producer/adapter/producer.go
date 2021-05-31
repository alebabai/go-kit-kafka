package adapter

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(producer sarama.SyncProducer) *Producer {
	return &Producer{
		producer: producer,
	}
}

func (p *Producer) Handle(ctx context.Context, msg *kafka.Message) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to produce message: %w", ctx.Err())
	default:
		if _, _, err := p.producer.SendMessage(TransformMessage(msg)); err != nil {
			return err
		}

		return nil
	}
}
