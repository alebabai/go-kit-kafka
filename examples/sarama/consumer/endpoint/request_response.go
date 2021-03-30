package endpoint

import (
	"github.com/alebabai/go-kit-kafka/examples/sarama/domain"
)

type CreateEventRequest struct {
	Payload *domain.Event
}

type CreateEventResponse struct {
}

type ListEventsRequest struct {
}

type ListEventsResponse struct {
	Results []domain.Event
}
