package adapter

import (
	kitkafka "github.com/alebabai/go-kit-kafka/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TransformMessage(msg *kafka.Message) *kitkafka.Message {
	headers := make([]kitkafka.Header, len(msg.Headers))
	for i, h := range msg.Headers {
		headers[i] = transformHeader(h)
	}

	return &kitkafka.Message{
		Topic:     *msg.TopicPartition.Topic,
		Partition: msg.TopicPartition.Partition,
		Offset:    int64(msg.TopicPartition.Offset),
		Key:       msg.Key,
		Value:     msg.Value,
		Headers:   headers,
		Timestamp: msg.Timestamp,
	}
}

func transformHeader(header kafka.Header) kitkafka.Header {
	return kitkafka.Header{
		Key:   []byte(header.Key),
		Value: header.Value,
	}
}
