package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alebabai/go-kafka"
	"github.com/alebabai/go-kit-kafka/v2/examples/common"
	"github.com/alebabai/go-kit-kafka/v2/transport"
	"github.com/go-kit/kit/endpoint"
)

func NewKafkaHandler(e endpoint.Endpoint) kafka.Handler {
	return transport.NewConsumer(
		e,
		decodeCreateEventKafkaRequest,
	)
}

func decodeCreateEventKafkaRequest(_ context.Context, msg *kafka.Message) (interface{}, error) {
	var out common.Event
	if err := json.Unmarshal(msg.Value, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create event request")
	}

	return CreateEventRequest{
		Payload: out,
	}, nil
}
