package kafka

import (
	"context"
)

type NoopErrorHandler struct {
}

func NewNoopErrorHandler() *NoopErrorHandler {
	return &NoopErrorHandler{}
}

func (eh *NoopErrorHandler) Handle(context.Context, error) {
}
