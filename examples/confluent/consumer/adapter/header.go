package adapter

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Header struct {
	h *kafka.Header
}

func NewHeader(h *kafka.Header) *Header {
	return &Header{
		h: h,
	}
}

func (h *Header) Key() string {
	return h.h.Key
}

func (h *Header) Value() []byte {
	return h.h.Value
}
