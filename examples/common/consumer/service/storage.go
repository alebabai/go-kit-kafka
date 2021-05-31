package service

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/go-kit/kit/log"

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

func (s *StorageService) Create(_ context.Context, e domain.Event) error {
	_ = s.logger.Log("msg", "saving an event", "event_id", e.ID)

	if _, ok := s.events[e.ID]; ok {
		return fmt.Errorf("event with id=%v already exists", e.ID)
	}

	s.m.Lock()
	s.events[e.ID] = e
	s.m.Unlock()

	return nil
}

func (s *StorageService) List(_ context.Context) ([]domain.Event, error) {
	s.m.Lock()

	out := make([]domain.Event, 0)
	for _, e := range s.events {
		out = append(out, e)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})

	// mark all viewed events as expired
	for k := range s.events {
		e := s.events[k]
		e.Expired = true
		s.events[k] = e
	}

	s.m.Unlock()

	return out, nil
}
