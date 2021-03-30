package adapter

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	kitkafka "github.com/alebabai/go-kit-kafka/kafka"
)

type Message struct {
	msg *kafka.Message
}

func NewMessage(msg *kafka.Message) *Message {
	return &Message{
		msg: msg,
	}
}

func (m *Message) Topic() string {
	return *m.msg.TopicPartition.Topic
}

func (m *Message) Partition() int32 {
	return m.msg.TopicPartition.Partition
}

func (m *Message) Offset() int64 {
	return int64(m.msg.TopicPartition.Offset)
}

func (m *Message) Key() []byte {
	return m.msg.Key
}

func (m *Message) Value() []byte {
	return m.msg.Value
}

func (m *Message) Headers() []kitkafka.Header {
	headers := make([]kitkafka.Header, 0)
	for _, h := range m.msg.Headers {
		headers = append(headers, NewHeader(&h))
	}

	return headers
}

func (m *Message) Timestamp() time.Time {
	return m.msg.Timestamp
}
