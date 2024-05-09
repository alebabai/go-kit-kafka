package consumer

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

func MakeCreateEventEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateEventRequest)

		res, err := svc.CreateEvent(ctx, req.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create event: %w", err)
		}

		return CreateEventResponse{
			Result: *res,
		}, nil
	}
}

func MakeListEventsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(ListEventsRequest)

		res, err := svc.ListEvents(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list events: %w", err)
		}

		return ListEventsResponse{
			Results: res,
		}, nil
	}
}
