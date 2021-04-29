package transport

import (
	"context"

	"github.com/alebabai/go-kit-kafka/kafka"
	"github.com/alebabai/go-kit-kafka/kafka/transport"

	"github.com/alebabai/go-kit-kafka/examples/common/producer/endpoint"
)

func NewKafkaProducer(handler kafka.Handler, topic string) *transport.Producer {
	return transport.NewProducer(handler, topic, encodeProduceEventRequest)
}

func encodeProduceEventRequest(ctx context.Context, msg *kafka.Message, request interface{}) error {
	req := request.(endpoint.ProduceEventRequest)
	return transport.EncodeJSONRequest(ctx, msg, req.Payload)
}
