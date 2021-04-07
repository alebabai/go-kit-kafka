package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/alebabai/go-kit-kafka/examples/sarama/consumer"
)

type Endpoints struct {
	CreateEventEndpoint endpoint.Endpoint
	ListEventsEndpoint  endpoint.Endpoint
}

func MakeCreateEventEndpoint(svc consumer.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateEventRequest)

		if err := svc.Create(ctx, *req.Payload); err != nil {
			return nil, fmt.Errorf("failed to create event: %w", err)
		}

		return CreateEventResponse{}, nil
	}
}

func MakeListEventsEndpoint(svc consumer.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(ListEventsRequest)

		ee, err := svc.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list events: %w", err)
		}

		return ListEventsResponse{
			Results: ee,
		}, nil
	}
}
