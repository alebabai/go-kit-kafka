package consumer

import (
	"context"

	"github.com/alebabai/go-kit-kafka/v2/examples/common/domain"
)

type Service interface {
	Create(ctx context.Context, e domain.Event) error
	List(ctx context.Context) ([]domain.Event, error)
}
