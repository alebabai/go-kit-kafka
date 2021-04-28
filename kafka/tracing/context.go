package tracing

import (
	"context"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type kafkaContextKey string

const (
	topicKey     kafkaContextKey = "kafka.topic"
	partitionKey kafkaContextKey = "kafka.partition"
	offsetKey    kafkaContextKey = "kafka.offset"
)

func MessageToContext(ctx context.Context, msg kafka.Message) context.Context {
	ctx = context.WithValue(ctx, topicKey, msg.Topic)
	ctx = context.WithValue(ctx, partitionKey, msg.Partition)
	ctx = context.WithValue(ctx, offsetKey, msg.Offset)

	return ctx
}

func ContextToTags(ctx context.Context) map[string]interface{} {
	tags := make(map[string]interface{})

	if topic := ctx.Value(topicKey); topic != nil {
		tags[string(topicKey)] = topic
	}

	if partition := ctx.Value(partitionKey); partition != nil {
		tags[string(partitionKey)] = partition
	}

	if offset := ctx.Value(offsetKey); offset != nil {
		tags[string(offsetKey)] = offset
	}

	return tags
}
