package consumer

import (
	"github.com/alebabai/go-kit-kafka/v2/examples/common/domain"
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
