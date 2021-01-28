package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/go-kit/kit/log"

	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
)

type producerMiddleware struct {
	topic    string
	producer *kafka.Producer
	logger   log.Logger
	next     Service
}

func (mv *producerMiddleware) Generate(ctx context.Context) (*domain.Event, error) {
	e, err := mv.next.Generate(ctx)
	if err != nil {
		return nil, err
	}

	if err := mv.produce(*e); err != nil {
		return nil, fmt.Errorf("failed to produce an event with ID=%s", e.ID)
	}

	return e, nil
}

func (mv *producerMiddleware) produce(e domain.Event) error {
	_ = mv.logger.Log("msg", "publishing an event", "event_id", e.ID)

	bytes, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal an event: %w", err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range mv.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					_ = mv.logger.Log("err", "failed to deliver message", "message_key", string(ev.Key), "message_payload", string(ev.Value))
				} else {
					_ = mv.logger.Log("msg", "message was successfully delivered", "message_key", string(ev.Key), "message_payload", string(ev.Value))
				}
			}
		}
	}()

	if err := mv.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &mv.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(e.ID),
		Value: bytes,
	}, nil); err != nil {
		return fmt.Errorf("failed to publish an event: %w", err)
	}

	return nil
}

func ProducerMiddleware(topic string, p *kafka.Producer, logger log.Logger) Middleware {
	return func(next Service) Service {
		return &producerMiddleware{
			topic:    topic,
			producer: p,
			logger:   logger,
			next:     next,
		}
	}
}
