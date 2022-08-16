package consumer

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(endpoints Endpoints) http.Handler {
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

	return r
}

func decodeListEventsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return ListEventsRequest{}, nil
}

func encodeListEventsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	httptransport.SetContentType("application/json")(ctx, w)

	res := response.(ListEventsResponse)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Results); err != nil {
		return err
	}

	return nil
}
