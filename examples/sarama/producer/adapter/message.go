package adapter

import (
	"github.com/Shopify/sarama"

	"github.com/alebabai/go-kit-kafka/kafka"
)

func TransformMessage(msg *kafka.Message) *sarama.ProducerMessage {
	headers := make([]sarama.RecordHeader, len(msg.Headers))
	for i, h := range msg.Headers {
		headers[i] = transformHeader(h)
	}

	return &sarama.ProducerMessage{
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Key:       sarama.ByteEncoder(msg.Key),
		Value:     sarama.ByteEncoder(msg.Value),
		Headers:   headers,
		Timestamp: msg.Timestamp,
	}
}

func transformHeader(header kafka.Header) sarama.RecordHeader {
	return sarama.RecordHeader{
		Key:   header.Key,
		Value: header.Value,
	}
}
