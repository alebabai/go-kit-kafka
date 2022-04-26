package producer

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GenerateEvent endpoint.Endpoint
}

func MakeGenerateEventEndpoint(svc Service) endpoint.Endpoint {
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
