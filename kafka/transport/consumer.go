package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type Consumer struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	before       []ConsumerRequestFunc
	after        []ConsumerResponseFunc
	finalizer    []ConsumerFinalizerFunc
	errorHandler transport.ErrorHandler
}

func NewConsumer(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	opts ...ConsumerOption,
) *Consumer {
	c := &Consumer{
		e:   e,
		dec: dec,
		errorHandler: transport.ErrorHandlerFunc(func(ctx context.Context, err error) {
			// noop
		}),
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Consumer) Handle(ctx context.Context, msg kafka.Message) (err error) {
	if len(c.finalizer) > 0 {
		defer func() {
			for _, f := range c.finalizer {
				f(ctx, msg, err)
			}
		}()
	}

	for _, f := range c.before {
		ctx = f(ctx, msg)
	}

	request, err := c.dec(ctx, msg)
	if err != nil {
		c.errorHandler.Handle(ctx, err)
		return err
	}

	response, err := c.e(ctx, request)
	if err != nil {
		c.errorHandler.Handle(ctx, err)
		return err
	}

	for _, f := range c.after {
		ctx = f(ctx, response)
	}

	return nil
}
