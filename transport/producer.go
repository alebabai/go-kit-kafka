package transport

import (
	"context"
	"encoding/json"

	"github.com/alebabai/go-kafka"
	"github.com/go-kit/kit/endpoint"
)

// Producer wraps single Kafka topic for message producing
// and implements endpoint.Endpoint.
type Producer struct {
	handler   kafka.Handler
	topic     string
	response  interface{}
	enc       EncodeRequestFunc
	before    []RequestFunc
	after     []ProducerResponseFunc
	finalizer []ProducerFinalizerFunc
}

// NewProducer constructs a new producer for a single Kafka topic.
func NewProducer(
	handler kafka.Handler,
	topic string,
	enc EncodeRequestFunc,
	options ...ProducerOption,
) *Producer {
	p := &Producer{
		handler:  handler,
		topic:    topic,
		response: struct{}{},
		enc:      enc,
	}
	for _, opt := range options {
		opt(p)
	}

	return p
}

// ProducerOption sets an optional parameter for a [Producer].
type ProducerOption func(producer *Producer)

// ProducerResponse sets the successful response value for a [Producer].
func ProducerResponse(response interface{}) ProducerOption {
	return func(p *Producer) {
		p.response = response
	}
}

// ProducerBefore sets the RequestFuncs that are applied to the outgoing producer
// request before it's invoked.
func ProducerBefore(before ...RequestFunc) ProducerOption {
	return func(p *Producer) {
		p.before = append(p.before, before...)
	}
}

// ProducerAfter adds one or more ProducerResponseFuncs, which are applied to the
// context after successful message producing.
// This is useful for context-manipulation operations.
func ProducerAfter(after ...ProducerResponseFunc) ProducerOption {
	return func(p *Producer) {
		p.after = append(p.after, after...)
	}
}

// ProducerFinalizer adds one or more ProducerFinalizerFuncs to be executed at the
// end of producing Kafka message. Finalizers are executed in the order in which they
// were added. By default, no finalizer is registered.
func ProducerFinalizer(f ...ProducerFinalizerFunc) ProducerOption {
	return func(p *Producer) {
		p.finalizer = append(p.finalizer, f...)
	}
}

// Endpoint returns a usable endpoint that invokes message producing.
func (p Producer) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if len(p.finalizer) > 0 {
			defer func() {
				for _, f := range p.finalizer {
					f(ctx, err)
				}
			}()
		}

		msg := &kafka.Message{
			Topic: p.topic,
		}
		if err := p.enc(ctx, msg, request); err != nil {
			return nil, err
		}

		for _, f := range p.before {
			ctx = f(ctx, msg)
		}

		if err := p.handler.Handle(ctx, *msg); err != nil {
			return nil, err
		}

		for _, f := range p.after {
			ctx = f(ctx)
		}

		return p.response, nil
	}
}

// ProducerFinalizerFunc can be used to perform work at the end of a producing Kafka message,
// after response is returned. The principal intended use is for error logging.
type ProducerFinalizerFunc func(ctx context.Context, err error)

// EncodeJSONRequest is an [EncodeRequestFunc] that serializes the request as a
// JSON object to the Message value.
// Many services can use it as a sensible default.
func EncodeJSONRequest(_ context.Context, msg *kafka.Message, request interface{}) error {
	rawJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg.Value = rawJSON

	return nil
}
