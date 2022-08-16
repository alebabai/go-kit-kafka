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
			decodeGenerateEventHTTPRequest,
			encodeGenerateEventHTTPResponse,
		))

	return router
}

func decodeGenerateEventHTTPRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return GenerateEventRequest{}, nil
}

func encodeGenerateEventHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	httptransport.SetContentType("application/json")(ctx, w)

	res := response.(GenerateEventResponse)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Result); err != nil {
		return fmt.Errorf("failed to encode json response: %w", err)
	}

	return nil
}
