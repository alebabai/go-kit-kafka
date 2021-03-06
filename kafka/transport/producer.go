package transport

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type Producer struct {
	handler   kafka.Handler
	topic     string
	response  interface{}
	enc       EncodeRequestFunc
	before    []RequestFunc
	after     []ProducerResponseFunc
	finalizer []ProducerFinalizerFunc
}

type successResponse struct{}

func NewProducer(
	handler kafka.Handler,
	topic string,
	enc EncodeRequestFunc,
	options ...ProducerOption,
) *Producer {
	p := &Producer{
		handler:  handler,
		topic:    topic,
		response: successResponse{},
		enc:      enc,
	}
	for _, opt := range options {
		opt(p)
	}

	return p
}

type ProducerOption func(consumer *Producer)

func ProducerResponse(response interface{}) ProducerOption {
	return func(p *Producer) {
		p.response = response
	}
}

func ProducerBefore(before ...RequestFunc) ProducerOption {
	return func(p *Producer) {
		p.before = append(p.before, before...)
	}
}

func ProducerAfter(after ...ProducerResponseFunc) ProducerOption {
	return func(p *Producer) {
		p.after = append(p.after, after...)
	}
}

func ProducerFinalizer(f ...ProducerFinalizerFunc) ProducerOption {
	return func(p *Producer) {
		p.finalizer = append(p.finalizer, f...)
	}
}

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

		if err := p.handler.Handle(ctx, msg); err != nil {
			return nil, err
		}

		for _, f := range p.after {
			ctx = f(ctx)
		}

		return p.response, nil
	}
}

type ProducerFinalizerFunc func(ctx context.Context, err error)

func EncodeJSONRequest(_ context.Context, msg *kafka.Message, request interface{}) error {
	rawJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg.Value = rawJSON

	return nil
}
