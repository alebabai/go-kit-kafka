package endpoint

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
)

type kafkaProducer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	Events() chan kafka.Event
}

func ProducerMiddleware(topic string, producer kafkaProducer, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			req := request.(GenerateEventRequest)

			response, err := next(ctx, req)
			if err != nil {
				return nil, err
			}

			resp := response.(GenerateEventResponse)
			if resp.Result != nil {
				if err := produceMessage(topic, producer, *resp.Result, logger); err != nil {
					err = fmt.Errorf("failed to produce message: %w", err)
					_ = level.Error(logger).Log("err", err)
				}
			}

			return response, nil
		}
	}
}

func produceMessage(topic string, p kafkaProducer, e domain.Event, logger log.Logger) error {
	_ = logger.Log("msg", "producing a message", "message_key", e.ID)

	bytes, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal an event: %w", err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					_ = level.Error(logger).Log(
						"err", "failed to deliver message",
						"message_key", string(ev.Key),
						"message_payload", string(ev.Value),
					)
				} else {
					_ = logger.Log(
						"msg", "message was successfully delivered",
						"message_key", string(ev.Key),
						"message_payload", string(ev.Value),
					)
				}
			}
		}
	}()

	if err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(e.ID),
		Value: bytes,
	}, nil); err != nil {
		return fmt.Errorf("failed to produce a message: %w", err)
	}

	return nil
}
