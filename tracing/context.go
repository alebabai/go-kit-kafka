package tracing

import (
	"context"

	"github.com/alebabai/go-kafka"
)

type KafkaContextKey string

const (
	KeyTopic     KafkaContextKey = "kafka.topic"
	KeyPartition KafkaContextKey = "kafka.partition"
	KeyOffset    KafkaContextKey = "kafka.offset"
)

// MessageToContext returns new context with topic, partition and offset values from msg.
func MessageToContext(ctx context.Context, msg *kafka.Message) context.Context {
	ctx = context.WithValue(ctx, KeyTopic, msg.Topic)
	ctx = context.WithValue(ctx, KeyPartition, msg.Partition)
	ctx = context.WithValue(ctx, KeyOffset, msg.Offset)

	return ctx
}

// ContextToTags returns new map of tags from ctx.
func ContextToTags(ctx context.Context) map[string]interface{} {
	tags := make(map[string]interface{})

	if topic := ctx.Value(KeyTopic); topic != nil {
		tags[string(KeyTopic)] = topic
	}

	if partition := ctx.Value(KeyPartition); partition != nil {
		tags[string(KeyPartition)] = partition
	}

	if offset := ctx.Value(KeyOffset); offset != nil {
		tags[string(KeyOffset)] = offset
	}

	return tags
}
