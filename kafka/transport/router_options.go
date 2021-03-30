package transport

import (
	"github.com/alebabai/go-kit-kafka/kafka"
)

type RouterOption func(*Router)

func RouterWithHandler(topic string, handler kafka.Handler) RouterOption {
	return func(r *Router) {
		r.AddHandler(topic, handler)
	}
}
