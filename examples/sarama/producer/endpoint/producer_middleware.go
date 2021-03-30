package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/alebabai/go-kit-kafka/examples/sarama/domain"
)

func ProducerMiddleware(topic string, producer sarama.SyncProducer, logger log.Logger) endpoint.Middleware {
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

func produceMessage(topic string, producer sarama.SyncProducer, e domain.Event, logger log.Logger) error {
	_ = logger.Log("msg", "producing a message", "message_key", e.ID)

	bytes, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal an event: %w", err)
	}

	if _, _, err := producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(e.ID),
			Value: sarama.ByteEncoder(bytes),
		}); err != nil {
		return fmt.Errorf("failed to produce a message: %w", err)
	}

	return nil
}
