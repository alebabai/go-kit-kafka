package producer

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

func Middleware(producer endpoint.Endpoint) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			response, err := next(ctx, request)
			if err != nil {
				return nil, err
			}

			if resp, ok := response.(GenerateEventResponse); ok {
				req := ProduceEventRequest{
					Payload: resp.Result,
				}
				if _, err = producer(ctx, req); err != nil {
					return nil, fmt.Errorf("failed to execute producer endpoint: %w", err)
				}
			}

			return response, err
		}
	}
}
