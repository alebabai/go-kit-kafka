package transport

import (
	"context"
	"github.com/alebabai/go-kit-kafka/examples/confluent/consumer/endpoint"
	"github.com/gorilla/mux"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(endpoints *endpoint.Endpoints) (http.Handler, error) {
	r := mux.
		NewRouter().
		StrictSlash(true)

	r.
		Path("/events").
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.ListEventsEndpoint,
			decodeListEventsRequest,
			encodeListEventsResponse,
		))

	return r, nil
}

func decodeListEventsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return endpoint.ListEventsRequest{}, nil
}

func encodeListEventsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoint.ListEventsResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Results); err != nil {
		return err
	}
	return nil
}