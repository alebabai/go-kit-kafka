package adapter

import (
	"github.com/Shopify/sarama"
)

type Header struct {
	h *sarama.RecordHeader
}

func NewHeader(h *sarama.RecordHeader) *Header {
	return &Header{
		h: h,
	}
}

func (h *Header) Key() []byte {
	return h.h.Key
}

func (h *Header) Value() []byte {
	return h.h.Value
}
