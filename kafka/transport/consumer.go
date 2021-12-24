package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"

	"github.com/alebabai/go-kit-kafka/kafka"
)

// Consumer wraps an endpoint and implements kafka.Handler.
type Consumer struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	before       []RequestFunc
	after        []ConsumerResponseFunc
	finalizer    []ConsumerFinalizerFunc
	errorHandler transport.ErrorHandler
}

// NewConsumer constructs a new consumer, which implements kafka.Handler and wraps
// the provided endpoint.
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

// ConsumerOption sets an optional parameter for consumer.
type ConsumerOption func(consumer *Consumer)

// ConsumerBefore functions are executed on the consumer message object
// before the request is decoded.
func ConsumerBefore(before ...RequestFunc) ConsumerOption {
	return func(c *Consumer) {
		c.before = append(c.before, before...)
	}
}

// ConsumerAfter functions are executed on the consumer reply after the
// endpoint is invoked, but before anything is published to the reply.
func ConsumerAfter(after ...ConsumerResponseFunc) ConsumerOption {
	return func(c *Consumer) {
		c.after = append(c.after, after...)
	}
}

// ConsumerErrorHandler is used to handle non-terminal errors. By default, non-terminal errors
// are ignored. This is intended as a diagnostic measure.
func ConsumerErrorHandler(errorHandler transport.ErrorHandler) ConsumerOption {
	return func(c *Consumer) {
		c.errorHandler = errorHandler
	}
}

// ConsumerFinalizer is executed at the end of every message processing.
// By default, no finalizer is registered.
func ConsumerFinalizer(f ...ConsumerFinalizerFunc) ConsumerOption {
	return func(c *Consumer) {
		c.finalizer = append(c.finalizer, f...)
	}
}

// Handle implements kafka.Handler.
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

// ConsumerFinalizerFunc can be used to perform work at the end of message processing,
// after the response has been constructed. The principal
// intended use is for request logging.
type ConsumerFinalizerFunc func(ctx context.Context, msg *kafka.Message, err error)
