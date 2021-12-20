package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/log"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type Consumer struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	before       []RequestFunc
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

type ConsumerOption func(consumer *Consumer)

func ConsumerBefore(before ...RequestFunc) ConsumerOption {
	return func(c *Consumer) {
		c.before = append(c.before, before...)
	}
}

func ConsumerAfter(after ...ConsumerResponseFunc) ConsumerOption {
	return func(c *Consumer) {
		c.after = append(c.after, after...)
	}
}

func ConsumerErrorLogger(logger log.Logger) ConsumerOption {
	return func(c *Consumer) {
		c.errorHandler = transport.NewLogErrorHandler(logger)
	}
}

func ConsumerErrorHandler(errorHandler transport.ErrorHandler) ConsumerOption {
	return func(c *Consumer) {
		c.errorHandler = errorHandler
	}
}

func ConsumerFinalizer(f ...ConsumerFinalizerFunc) ConsumerOption {
	return func(c *Consumer) {
		c.finalizer = append(c.finalizer, f...)
	}
}

func (c Consumer) Handle(ctx context.Context, msg *kafka.Message) (err error) {
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

type ConsumerFinalizerFunc func(ctx context.Context, msg *kafka.Message, err error)
