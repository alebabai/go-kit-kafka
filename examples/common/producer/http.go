package producer

import (
	"context"
	"fmt"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(endpoints Endpoints) http.Handler {
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
	return GenerateEventRequest{}, nil
}

func encodeGenerateEventResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(GenerateEventResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Result); err != nil {
		return fmt.Errorf("failed to encode json response: %w", err)
	}

	return nil
}
