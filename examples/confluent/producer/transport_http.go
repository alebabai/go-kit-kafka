package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(e *Endpoints) (http.Handler, error) {
	r := mux.
		NewRouter().
		StrictSlash(true)

	r.
		Path("/events").
		Methods("POST").
		Handler(httptransport.NewServer(
			e.GenerateEventEndpoint,
			decodeGenerateEventRequest,
			encodeGenerateEventResponse,
		))

	return r, nil
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
