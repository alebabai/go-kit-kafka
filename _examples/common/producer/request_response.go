package producer

import (
	"github.com/alebabai/go-kit-kafka/v2/examples/common"
)

type GenerateEventRequest struct {
}

type GenerateEventResponse struct {
	Result common.Event
}

type ProduceEventRequest struct {
	Payload common.Event
}

type ProduceEventResponse struct {
}
