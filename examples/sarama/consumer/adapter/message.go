package adapter

import (
	"github.com/Shopify/sarama"
	"time"

	kitkafka "github.com/alebabai/go-kit-kafka/kafka"
)

type Message struct {
	msg *sarama.ConsumerMessage
}

func NewMessage(msg *sarama.ConsumerMessage) *Message {
	return &Message{
		msg: msg,
	}
}

func (m *Message) Topic() string {
	return m.msg.Topic
}

func (m *Message) Partition() int32 {
	return m.msg.Partition
}

func (m *Message) Offset() int64 {
	return m.msg.Offset
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
		headers = append(headers, NewHeader(h))
	}

	return headers
}

func (m *Message) Timestamp() time.Time {
	return m.msg.Timestamp
}
