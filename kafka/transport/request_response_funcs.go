package transport

import (
	"context"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type ConsumerRequestFunc func(ctx context.Context, msg kafka.Message) context.Context

type ConsumerResponseFunc func(ctx context.Context, response interface{}) context.Context

type ConsumerFinalizerFunc func(ctx context.Context, msg kafka.Message, err error)
