package opentracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/go-kit/kit/endpoint"
)

func TraceConsumer(tracer opentracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			span := opentracing.SpanFromContext(ctx)
			if span == nil {
				span = tracer.StartSpan(operationName)
			} else {
				span.SetOperationName(operationName)
			}
			defer span.Finish()

			ext.SpanKindConsumer.Set(span)
			ctx = opentracing.ContextWithSpan(ctx, span)

			return next(ctx, request)
		}
	}
}
