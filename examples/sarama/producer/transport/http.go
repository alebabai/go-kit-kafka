package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/alebabai/go-kit-kafka/examples/sarama/producer/endpoint"
)

func NewHTTPHandler(endpoints endpoint.Endpoints) http.Handler {
	router := mux.
		NewRouter().
		StrictSlash(true)

	router.
		Path("/events").
		Methods("POST").
		Handler(httptransport.NewServer(
			endpoints.GenerateEvent,
			decodeGenerateEventRequest,
			encodeGenerateEventResponse,
		))

	return router
}

func decodeGenerateEventRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return endpoint.GenerateEventRequest{}, nil
}

func encodeGenerateEventResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoint.GenerateEventResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Result); err != nil {
		return fmt.Errorf("failed to encode json response: %w", err)
	}

	return nil
}
