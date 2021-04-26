package transport

import (
	"context"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type RequestFunc func(ctx context.Context, msg kafka.Message) context.Context

type ConsumerResponseFunc func(ctx context.Context, response interface{}) context.Context
