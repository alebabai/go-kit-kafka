package adapter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"

	kitkafka "github.com/alebabai/go-kit-kafka/kafka"
)

type consumer interface {
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	io.Closer
}

type Listener struct {
	consumer consumer
	handler  kitkafka.Handler

	readTimeout time.Duration

	errorHandler transport.ErrorHandler
}

func NewListener(consumer consumer, handler kitkafka.Handler, opts ...ListenerOption) (*Listener, error) {
	if consumer == nil {
		return nil, errors.New("consumer cannot be nil")
	}
	if handler == nil {
		return nil, errors.New("handler cannot be nil")
	}

	l := &Listener{
		consumer: consumer,
		handler:  handler,

		readTimeout: -1,
		errorHandler: transport.ErrorHandlerFunc(func(context.Context, error) {
			// noop
		}),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

type ListenerOption func(*Listener)

func ListenerErrorLogger(logger log.Logger) ListenerOption {
	return func(l *Listener) {
		l.errorHandler = transport.NewLogErrorHandler(logger)
	}
}

func ListenerErrorHandler(errHandler transport.ErrorHandler) ListenerOption {
	return func(l *Listener) {
		l.errorHandler = errHandler
	}
}

func ListenerReadTimeout(readTimeout time.Duration) ListenerOption {
	return func(l *Listener) {
		l.readTimeout = readTimeout
	}
}

func (l *Listener) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := l.consumer.ReadMessage(l.readTimeout)
			if err != nil {
				err = fmt.Errorf("failed to read kafka message: %w", err)
				l.errorHandler.Handle(ctx, err)
				continue
			}

			if err := l.handler.Handle(ctx, *TransformMessage(*msg)); err != nil {
				err = fmt.Errorf("failed to handle kafka message: %w", err)
				l.errorHandler.Handle(ctx, err)
				continue
			}
		}
	}
}
