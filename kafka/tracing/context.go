package tracing

import (
	"context"

	"github.com/alebabai/go-kit-kafka/kafka"
)

type kafkaContextKey string

const (
	keyTopic     kafkaContextKey = "kafka.topic"
	keyPartition kafkaContextKey = "kafka.partition"
	keyOffset    kafkaContextKey = "kafka.offset"
)

func MessageToContext(ctx context.Context, msg *kafka.Message) context.Context {
	ctx = context.WithValue(ctx, keyTopic, msg.Topic)
	ctx = context.WithValue(ctx, keyPartition, msg.Partition)
	ctx = context.WithValue(ctx, keyOffset, msg.Offset)

	return ctx
}

func ContextToTags(ctx context.Context) map[string]interface{} {
	tags := make(map[string]interface{})

	if topic := ctx.Value(keyTopic); topic != nil {
		tags[string(keyTopic)] = topic
	}

	if partition := ctx.Value(keyPartition); partition != nil {
		tags[string(keyPartition)] = partition
	}

	if offset := ctx.Value(keyOffset); offset != nil {
		tags[string(keyOffset)] = offset
	}

	return tags
}
