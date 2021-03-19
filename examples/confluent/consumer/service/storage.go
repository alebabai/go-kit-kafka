package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kit/kit/log"

	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
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
	_ = s.logger.Log("msg", "saving an domain.Event", "event_id", e.ID)

	if _, ok := s.events[e.ID]; ok {
		return fmt.Errorf("domain.Event with id=%v already exists", e.ID)
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

	// mark all viewed events as expired
	for k := range s.events {
		e := s.events[k]
		e.Expired = true
		s.events[k] = e
	}

	s.m.Unlock()

	return out, nil
}
