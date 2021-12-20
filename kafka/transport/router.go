package transport

import (
	"context"
	"fmt"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type Router map[string][]kafka.Handler

func (r Router) AddHandler(topic string, handler kafka.Handler) Router {
	r[topic] = append(r[topic], handler)

	return r
}

func (r Router) Handle(ctx context.Context, msg *kafka.Message) error {
	for _, h := range r[msg.Topic] {
		if err := h.Handle(ctx, msg); err != nil {
			return fmt.Errorf("failed to handle message from kafka topic=%s: %w", msg.Topic, err)
		}
	}

	return nil
}
