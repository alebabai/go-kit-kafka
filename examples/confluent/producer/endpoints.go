package main

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GenerateEventEndpoint endpoint.Endpoint
}

func NewEndpoints(svc Service) (*Endpoints, error) {
	return &Endpoints{
		GenerateEventEndpoint: makeGenerateEventEndpoint(svc),
	}, nil
}

func makeGenerateEventEndpoint(svc Service) endpoint.Endpoint {
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
