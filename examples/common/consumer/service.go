package consumer

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/alebabai/go-kit-kafka/v2/examples/common"
	"github.com/go-kit/log"
)

type Service interface {
	CreateEvent(ctx context.Context, in common.Event) (*common.Event, error)
	ListEvents(ctx context.Context) ([]common.Event, error)
}

type service struct {
	cache  map[string]common.Event
	m      sync.Mutex
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	return &service{
		cache:  make(map[string]common.Event),
		logger: logger,
	}
}

func (svc *service) CreateEvent(_ context.Context, in common.Event) (*common.Event, error) {
	_ = svc.logger.Log("msg", "saving an event", "event_id", in.ID)

	if _, ok := svc.cache[in.ID]; ok {
		return nil, fmt.Errorf("event with id=%v already exists", in.ID)
	}

	in.State = "new"

	svc.m.Lock()
	svc.cache[in.ID] = in
	svc.m.Unlock()

	return &in, nil
}

func (svc *service) ListEvents(_ context.Context) ([]common.Event, error) {
	svc.m.Lock()

	out := make([]common.Event, 0)
	for _, e := range svc.cache {
		out = append(out, e)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})

	// mark all viewed events as expired
	for k := range svc.cache {
		e := svc.cache[k]
		e.State = "expired"
		svc.cache[k] = e
	}

	svc.m.Unlock()

	return out, nil
}
