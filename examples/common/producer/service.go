package producer

import (
	"context"
	"time"

	"github.com/alebabai/go-kit-kafka/v2/examples/common"
	"github.com/go-kit/log"
	"github.com/google/uuid"
)

type Service interface {
	GenerateEvent(ctx context.Context) (*common.Event, error)
}

type service struct {
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (svc *service) GenerateEvent(_ context.Context) (*common.Event, error) {
	_ = svc.logger.Log("msg", "generating an event")

	return &common.Event{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		State:     "generated",
	}, nil
}
