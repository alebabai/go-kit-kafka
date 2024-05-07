package producer

import (
	"context"

	"github.com/alebabai/go-kafka"
	"github.com/alebabai/go-kit-kafka/v2/tracing"
	"github.com/alebabai/go-kit-kafka/v2/transport"
)

func NewKafkaProducer(handler kafka.Handler, topic string) *transport.Producer {
	return transport.NewProducer(
		handler,
		topic,
		encodeProduceEventKafkaRequest,
		transport.ProducerBefore(tracing.MessageToContext),
	)
}

func encodeProduceEventKafkaRequest(ctx context.Context, msg *kafka.Message, request interface{}) error {
	req := request.(ProduceEventRequest)

	return transport.EncodeJSONRequest(ctx, msg, req.Payload)
}
