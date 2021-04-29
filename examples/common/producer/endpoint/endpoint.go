package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/alebabai/go-kit-kafka/examples/common/producer"
)

type Endpoints struct {
	GenerateEvent endpoint.Endpoint
}

func MakeGenerateEventEndpoint(svc producer.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(GenerateEventRequest)

		t, err := svc.Generate(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate an event: %w", err)
		}

		return GenerateEventResponse{
			Result: t,
		}, nil
	}
}
