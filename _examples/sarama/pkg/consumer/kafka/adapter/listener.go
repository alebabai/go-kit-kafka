package adapter

import (
	"context"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/transport"
)

type Listener struct {
	topics []string

	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler sarama.ConsumerGroupHandler

	errorHandler transport.ErrorHandler
}

type ListenerOption func(*Listener)

func NewListener(
	topics []string,
	consumerGroup sarama.ConsumerGroup,
	consumerGroupHandler sarama.ConsumerGroupHandler,
	opts ...ListenerOption,
) (*Listener, error) {
	if len(topics) == 0 {
		return nil, errors.New("topics cannot be empty")
	}

	if consumerGroup == nil {
		return nil, errors.New("consumer group cannot be nil")
	}

	if consumerGroupHandler == nil {
		return nil, errors.New("consumer group handler cannot be nil")
	}

	l := &Listener{
		topics: topics,

		consumerGroup:        consumerGroup,
		consumerGroupHandler: consumerGroupHandler,

		errorHandler: transport.ErrorHandlerFunc(func(context.Context, error) {
			// noop
		}),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

func ListenerErrorHandler(errHandler transport.ErrorHandler) ListenerOption {
	return func(l *Listener) {
		l.errorHandler = errHandler
	}
}

func (l *Listener) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := l.consumerGroup.Consume(ctx, l.topics, l.consumerGroupHandler); err != nil {
				l.errorHandler.Handle(ctx, err)
				return err
			}
		}
	}
}
