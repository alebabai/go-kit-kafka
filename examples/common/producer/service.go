package producer

import (
	"context"

	"github.com/alebabai/go-kit-kafka/examples/common/domain"
)

type Service interface {
	Generate(ctx context.Context) (*domain.Event, error)
}
