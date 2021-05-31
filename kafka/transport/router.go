package transport

import (
	"context"
	"fmt"

	"github.com/alebabai/go-kit-kafka/kafka"
)

// Handlers represents Topic -> []Handler mapping
type Handlers map[string][]kafka.Handler

type Router struct {
	handlers Handlers
}

func NewRouter(opts ...RouterOption) *Router {
	r := &Router{
		handlers: make(Handlers),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type RouterOption func(*Router)

func RouterWithHandler(topic string, handler kafka.Handler) RouterOption {
	return func(r *Router) {
		r.AddHandler(topic, handler)
	}
}

func RouterWithHandlers(handlers Handlers) RouterOption {
	return func(r *Router) {
		r.handlers = handlers
	}
}

func (r *Router) AddHandler(topic string, handler kafka.Handler) *Router {
	if len(r.handlers) == 0 {
		r.handlers = make(Handlers)
	}

	r.handlers[topic] = append(r.handlers[topic], handler)

	return r
}

func (r Router) Handlers() Handlers {
	return r.handlers
}

func (r Router) Handle(ctx context.Context, msg *kafka.Message) error {
	for _, h := range r.handlers[msg.Topic] {
		if err := h.Handle(ctx, msg); err != nil {
			return fmt.Errorf("failed to handle Message from kafka topic=%s: %w", msg.Topic, err)
		}
	}

	return nil
}
