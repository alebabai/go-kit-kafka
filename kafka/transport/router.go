package transport

import (
	"context"
	"errors"
	"fmt"

	"github.com/alebabai/go-kit-kafka/kafka"
)

// Handlers represents Topic -> Handler mapping
type Handlers map[string]kafka.Handler

type Router struct {
	handlers Handlers
}

func NewRouter(handlers Handlers, opts ...RouterOption) (*Router, error) {
	if len(handlers) == 0 {
		return nil, errors.New("handlers cannot be empty")
	}

	l := &Router{
		handlers: handlers,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

func (r *Router) Handle(ctx context.Context, msg kafka.Message) error {
	h := r.handlers[msg.Topic()]
	if h != nil {
		if err := h.Handle(ctx, msg); err != nil {
			return fmt.Errorf("failed to handle kafka message from topic=%s: %w", msg.Topic(), err)
		}
	}

	return nil
}
