package transport

import (
	"context"

	"github.com/alebabai/go-kit-kafka/kafka"
)

// DecodeRequestFunc extracts a user-domain request object from
// an Kafka message. It is designed to be used in Kafka Consumers.
type DecodeRequestFunc func(ctx context.Context, msg *kafka.Message) (request interface{}, err error)

// EncodeRequestFunc encodes the passed request object into
// an Kafka message object. It is designed to be used in Kafka Producers.
type EncodeRequestFunc func(context.Context, *kafka.Message, interface{}) error

// EncodeResponseFunc encodes the passed response object into
// an Kafka message object. It is designed to be used in Kafka Consumers.
type EncodeResponseFunc func(context.Context, interface{}) error
