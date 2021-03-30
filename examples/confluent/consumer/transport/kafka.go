package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alebabai/go-kit-kafka/kafka"
	"github.com/alebabai/go-kit-kafka/kafka/transport"

	"github.com/alebabai/go-kit-kafka/examples/confluent/consumer/endpoint"
	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
)

func NewKafkaHandler(e *endpoint.Endpoints) kafka.Handler {
	return transport.NewConsumer(e.CreateEventEndpoint, decodeCreateEventRequest)
}

func decodeCreateEventRequest(ctx context.Context, msg kafka.Message) (interface{}, error) {
	var e domain.Event
	if err := json.Unmarshal(msg.Value(), &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create event request")
	}

	return endpoint.CreateEventRequest{
		Payload: &e,
	}, nil
}
