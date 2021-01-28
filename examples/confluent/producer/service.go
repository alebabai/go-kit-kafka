package main

import (
	"context"

	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
)

type Service interface {
	Generate(ctx context.Context) (*domain.Event, error)
}

type Middleware func(svc Service) Service
