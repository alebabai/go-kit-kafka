package opentracing

import (
	"github.com/go-kit/kit/endpoint"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// TraceConsumer returns a [endpoint.Middleware] that wraps the next [endpoint.Endpoint] in an
// [opentracing.Span] called operationName with consumer span.kind tag.
func TraceConsumer(tracer opentracing.Tracer, operationName string, opts ...kitopentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, kitopentracing.WithTags(map[string]interface{}{
		ext.SpanKindConsumer.Key: ext.SpanKindConsumer.Value,
	}))

	return kitopentracing.TraceEndpoint(tracer, operationName, opts...)
}

// TraceProducer returns a [endpoint.Middleware] that wraps the next [endpoint.Endpoint] in an
// [opentracing.Span] called operationName with producer span.kind tag.
func TraceProducer(tracer opentracing.Tracer, operationName string, opts ...kitopentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, kitopentracing.WithTags(map[string]interface{}{
		ext.SpanKindProducer.Key: ext.SpanKindProducer.Value,
	}))

	return kitopentracing.TraceEndpoint(tracer, operationName, opts...)
}
