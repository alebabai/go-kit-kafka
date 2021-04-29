package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/Shopify/sarama"

	"github.com/go-kit/kit/transport"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type ConsumerGroupHandler struct {
	handler      kafka.Handler
	errorHandler transport.ErrorHandler
}

func NewConsumerGroupHandler(handler kafka.Handler, opts ...ConsumerGroupHandlerOption) (*ConsumerGroupHandler, error) {
	if handler == nil {
		return nil, errors.New("handler cannot be nil")
	}

	l := &ConsumerGroupHandler{
		handler: handler,
		errorHandler: transport.ErrorHandlerFunc(func(context.Context, error) {
			// noop
		}),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()
	for msg := range claim.Messages() {
		if err := h.handler.Handle(ctx, *TransformMessage(*msg)); err != nil {
			err = fmt.Errorf("failed to handle kafka message: %w", err)
			h.errorHandler.Handle(ctx, err)
			continue
		}
		session.MarkMessage(msg, "")
	}

	return nil
}
