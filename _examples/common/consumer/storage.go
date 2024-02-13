package consumer

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/go-kit/log"

	"github.com/alebabai/go-kit-kafka/examples/common/domain"
)

type StorageService struct {
	events map[string]domain.Event

	m sync.Mutex

	logger log.Logger
}

func NewStorageService(logger log.Logger) (*StorageService, error) {
	return &StorageService{
		events: make(map[string]domain.Event),
		logger: logger,
	}, nil
}

func (svc *StorageService) Create(_ context.Context, e domain.Event) error {
	_ = svc.logger.Log("msg", "saving an event", "event_id", e.ID)

	if _, ok := svc.events[e.ID]; ok {
		return fmt.Errorf("event with id=%v already exists", e.ID)
	}

	svc.m.Lock()
	svc.events[e.ID] = e
	svc.m.Unlock()

	return nil
}

func (svc *StorageService) List(_ context.Context) ([]domain.Event, error) {
	svc.m.Lock()

	out := make([]domain.Event, 0)
	for _, e := range svc.events {
		out = append(out, e)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})

	// mark all viewed events as expired
	for k := range svc.events {
		e := svc.events[k]
		e.Expired = true
		svc.events[k] = e
	}

	svc.m.Unlock()

	return out, nil
}
