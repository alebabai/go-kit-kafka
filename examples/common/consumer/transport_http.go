package consumer

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(e endpoint.Endpoint) http.Handler {
	m := http.NewServeMux()
	m.Handle("/events", httptransport.NewServer(
		e,
		decodeListEventsHTTPRequest,
		encodeListEventsHTTPResponse,
	))

	return m
}

func decodeListEventsHTTPRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return ListEventsRequest{}, nil
}

func encodeListEventsHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	httptransport.SetContentType("application/json")(ctx, w)

	resp := response.(ListEventsResponse)
	if err := httptransport.EncodeJSONResponse(ctx, w, resp.Results); err != nil {
		return err
	}

	return nil
}
