package transport

import (
	"context"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type DecodeRequestFunc func(ctx context.Context, msg kafka.Message) (req interface{}, err error)

type EncodeRequestFunc func(context.Context, *kafka.Message, interface{}) error

type EncodeResponseFunc func(context.Context, interface{}) error
