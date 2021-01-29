package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/transport"
)

// Handlers represents Topic -> Handler mapping
type Handlers map[string]Handler

type Listener struct {
	reader   Reader
	handlers Handlers

	manualCommit bool

	errorHandler transport.ErrorHandler
}

func NewListener(reader Reader, handlers Handlers, opts ...ListenerOption) (*Listener, error) {
	if reader == nil {
		return nil, errors.New("reader cannot be nil")
	}
	if len(handlers) == 0 {
		return nil, errors.New("handlers cannot be empty")
	}

	l := &Listener{
		reader:   reader,
		handlers: handlers,

		errorHandler: NewNoopErrorHandler(),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

func (l *Listener) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := l.reader.ReadMessage(ctx)
			if err != nil {
				err = fmt.Errorf("failed to read kafka message: %w", err)
				l.errorHandler.Handle(ctx, err)
				continue
			}

			l.onMessage(ctx, msg)
		}
	}
}

func (l *Listener) onMessage(ctx context.Context, msg Message) {
	h := l.handlers[msg.Topic()]
	if h != nil {
		if err := h.Handle(ctx, msg); err != nil {
			err = fmt.Errorf("failed to handle kafka message from topic=%s: %w", msg.Topic(), err)
			l.errorHandler.Handle(ctx, err)
		} else if l.manualCommit {
			if err := l.reader.Commit(ctx); err != nil {
				err = fmt.Errorf("failed to commit kafka offsets: %w", err)
				l.errorHandler.Handle(ctx, err)
			}
		}
	}
}

func (l *Listener) Close() error {
	if l.reader != nil {
		return l.reader.Close()
	}

	return nil
}