package producer

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/google/uuid"

	"github.com/alebabai/go-kit-kafka/examples/common/domain"
)

type GeneratorService struct {
	logger log.Logger
}

func NewGeneratorService(logger log.Logger) (*GeneratorService, error) {
	return &GeneratorService{
		logger: logger,
	}, nil
}

func (svc *GeneratorService) Generate(_ context.Context) (*domain.Event, error) {
	_ = svc.logger.Log("msg", "generating an event")

	return &domain.Event{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Expired:   false,
	}, nil
}
