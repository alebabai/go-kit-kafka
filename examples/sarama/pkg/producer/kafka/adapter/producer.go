package adapter

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/alebabai/go-kafka"
	adapter "github.com/alebabai/go-kafka/adapter/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(producer sarama.SyncProducer) *Producer {
	return &Producer{
		producer: producer,
	}
}

func (p *Producer) Handle(ctx context.Context, msg kafka.Message) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to produce message: %w", ctx.Err())
	default:
		pmsg := adapter.ConvertKafkaMessageToProducerMessage(msg)
		if _, _, err := p.producer.SendMessage(&pmsg); err != nil {
			return err
		}

		return nil
	}
}
