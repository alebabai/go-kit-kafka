package opentracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/go-kit/kit/endpoint"
)

func TraceEndpoint(tracer opentracing.Tracer, operationName string, opts ...EndpointOption) endpoint.Middleware {
	cfg := &EndpointOptions{
		Tags: make(opentracing.Tags),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if cfg.GetOperationName != nil {
				if newOperationName := cfg.GetOperationName(ctx, operationName); newOperationName != "" {
					operationName = newOperationName
				}
			}

			var span opentracing.Span
			if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
				span = tracer.StartSpan(
					operationName,
					opentracing.ChildOf(parentSpan.Context()),
				)
			} else {
				span = tracer.StartSpan(operationName)
			}
			defer span.Finish()

			applyTags(span, cfg.Tags)
			if cfg.GetTags != nil {
				extraTags := cfg.GetTags(ctx)
				applyTags(span, extraTags)
			}

			ctx = opentracing.ContextWithSpan(ctx, span)

			response, err := next(ctx, request)
			if err := identifyError(response, err, cfg.IgnoreBusinessError); err != nil {
				ext.LogError(span, err)
			}

			return response, err
		}
	}
}

func TraceConsumer(tracer opentracing.Tracer, operationName string, opts ...EndpointOption) endpoint.Middleware {
	opts = append(
		opts,
		WithTags(map[string]interface{}{
			ext.SpanKindConsumer.Key: ext.SpanKindConsumer.Value,
		}),
	)

	return TraceEndpoint(tracer, operationName, opts...)
}

func applyTags(span opentracing.Span, tags opentracing.Tags) {
	for key, value := range tags {
		span.SetTag(key, value)
	}
}

func identifyError(response interface{}, err error, ignoreBusinessError bool) error {
	if err != nil {
		return err
	}

	if !ignoreBusinessError {
		if res, ok := response.(endpoint.Failer); ok {
			return res.Failed()
		}
	}

	return nil
}