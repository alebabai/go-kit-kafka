package transport

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/alebabai/go-kit-kafka/examples/confluent/consumer/endpoint"
)

func NewHTTPHandler(endpoints endpoint.Endpoints) http.Handler {
	router := mux.
		NewRouter().
		StrictSlash(true)

	router.
		Path("/events").
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.ListEventsEndpoint,
			decodeListEventsRequest,
			encodeListEventsResponse,
		))

	return router
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
