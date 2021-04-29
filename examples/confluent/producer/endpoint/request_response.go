package endpoint

import (
	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
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
