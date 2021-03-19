package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	kitendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/alebabai/go-kit-kafka/examples/confluent/producer/endpoint"
)

func NewHTTPHandler(e kitendpoint.Endpoint) (http.Handler, error) {
	r := mux.
		NewRouter().
		StrictSlash(true)

	r.
		Path("/events").
		Methods("POST").
		Handler(httptransport.NewServer(
			e,
			decodeGenerateEventRequest,
			encodeGenerateEventResponse,
		))

	return r, nil
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
