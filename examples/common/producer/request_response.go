package producer

import (
	"github.com/alebabai/go-kit-kafka/examples/common/domain"
)

type GenerateEventRequest struct {
}

type GenerateEventResponse struct {
	Result *domain.Event
}

type ProduceEventRequest struct {
	Payload *domain.Event
}

type ProduceEventResponse struct {
}
