package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alebabai/go-kafka"
	"github.com/alebabai/go-kit-kafka/v2/transport"
	"github.com/go-kit/kit/endpoint"

	"github.com/alebabai/go-kit-kafka/v2/examples/common/domain"
)

func NewKafkaHandler(e endpoint.Endpoint) kafka.Handler {
	return transport.NewConsumer(
		e,
		decodeCreateEventKafkaRequest,
	)
}

func decodeCreateEventKafkaRequest(_ context.Context, msg *kafka.Message) (interface{}, error) {
	var e domain.Event
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create event request")
	}

	return CreateEventRequest{
		Payload: &e,
	}, nil
}
