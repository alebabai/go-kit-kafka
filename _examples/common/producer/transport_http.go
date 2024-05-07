package producer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(e endpoint.Endpoint) http.Handler {
	m := http.NewServeMux()
	m.Handle("/events", httptransport.NewServer(
		e,
		decodeGenerateEventHTTPRequest,
		encodeGenerateEventHTTPResponse,
	))

	return m
}

func decodeGenerateEventHTTPRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return GenerateEventRequest{}, nil
}

func encodeGenerateEventHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	httptransport.SetContentType("application/json")(ctx, w)

	resp := response.(GenerateEventResponse)
	if err := httptransport.EncodeJSONResponse(ctx, w, resp.Result); err != nil {
		return fmt.Errorf("failed to encode json response: %w", err)
	}

	return nil
}
