package adapter

import (
	"github.com/Shopify/sarama"
	"github.com/alebabai/go-kafka"
)

func TransformMessage(msg *sarama.ConsumerMessage) *kafka.Message {
	headers := make([]kafka.Header, len(msg.Headers))
	for i, h := range msg.Headers {
		headers[i] = transformHeader(*h)
	}

	return &kafka.Message{
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Key:       msg.Key,
		Value:     msg.Value,
		Headers:   headers,
		Timestamp: msg.Timestamp,
	}
}

func transformHeader(header sarama.RecordHeader) kafka.Header {
	return kafka.Header{
		Key:   header.Key,
		Value: header.Value,
	}
}
