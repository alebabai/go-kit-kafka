package main

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/go-kit/kit/log"

	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
)

type Generator struct {
	logger log.Logger
}

func NewGenerator(logger log.Logger) (*Generator, error) {
	return &Generator{
		logger: logger,
	}, nil
}

func (t *Generator) Generate(ctx context.Context) (*domain.Event, error) {
	_ = t.logger.Log("msg", "generating an event")

	return &domain.Event{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Expired:   false,
	}, nil
}
