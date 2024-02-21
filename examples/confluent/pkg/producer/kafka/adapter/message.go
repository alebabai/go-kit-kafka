package adapter

import (
	kitkafka "github.com/alebabai/go-kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TransformMessage(msg *kitkafka.Message) *kafka.Message {
	headers := make([]kafka.Header, len(msg.Headers))
	for i, h := range msg.Headers {
		headers[i] = transformHeader(h)
	}

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &msg.Topic,
			Partition: msg.Partition,
			Offset:    kafka.Offset(msg.Offset),
		},
		Key:       msg.Key,
		Value:     msg.Value,
		Headers:   headers,
		Timestamp: msg.Timestamp,
	}
}

func transformHeader(header kitkafka.Header) kafka.Header {
	return kafka.Header{
		Key:   string(header.Key),
		Value: header.Value,
	}
}
