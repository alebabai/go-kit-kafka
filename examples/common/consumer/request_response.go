package consumer

import (
	"github.com/alebabai/go-kit-kafka/v2/examples/common"
)

type CreateEventRequest struct {
	Payload common.Event
}

type CreateEventResponse struct {
	Result common.Event
}

type ListEventsRequest struct {
}

type ListEventsResponse struct {
	Results []common.Event
}
