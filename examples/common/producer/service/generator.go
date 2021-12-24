package service

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/google/uuid"

	"github.com/alebabai/go-kit-kafka/examples/common/domain"
)

type Generator struct {
	logger log.Logger
}

func NewGeneratorService(logger log.Logger) (*Generator, error) {
	return &Generator{
		logger: logger,
	}, nil
}

func (t *Generator) Generate(_ context.Context) (*domain.Event, error) {
	_ = t.logger.Log("msg", "generating an event")

	return &domain.Event{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Expired:   false,
	}, nil
}
