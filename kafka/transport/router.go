package transport

import (
	"context"
	"fmt"

	"github.com/alebabai/go-kit-kafka/kafka"
)

// Router represents mapping topic -> []kafka.Handler
// and implements kafka.Handler with routing handlers by topic.
type Router map[string][]kafka.Handler

// AddHandler appends kafka.Handler for specific topic.
func (r Router) AddHandler(topic string, handler kafka.Handler) Router {
	r[topic] = append(r[topic], handler)

	return r
}

// Handle implements kafka.Handler.
func (r Router) Handle(ctx context.Context, msg *kafka.Message) error {
	for _, h := range r[msg.Topic] {
		if err := h.Handle(ctx, msg); err != nil {
			return fmt.Errorf("failed to handle message from kafka topic=%s: %w", msg.Topic, err)
		}
	}

	return nil
}
