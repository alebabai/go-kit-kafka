package transport

import (
	"context"

	"github.com/alebabai/go-kit-kafka/kafka"
)

// RequestFunc may take information from a Kafka message and put it into a
// request context. In Consumers, RequestFuncs are executed prior to invoking the
// endpoint.
type RequestFunc func(ctx context.Context, msg *kafka.Message) context.Context

// ConsumerResponseFunc may take information from a request context and use it to
// manipulate a Producer. ConsumerResponseFuncs are only executed in
// consumers, after invoking the endpoint but prior to publishing a reply.
type ConsumerResponseFunc func(ctx context.Context, response interface{}) context.Context

// ProducerResponseFunc may take information from a request context.
// ProducerResponseFunc are only executed in producers, after a request has been produced.
type ProducerResponseFunc func(ctx context.Context) context.Context
