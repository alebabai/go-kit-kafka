package opentracing

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"

	"github.com/alebabai/go-kit-kafka/kafka"
	"github.com/alebabai/go-kit-kafka/kafka/transport"
)

func TraceConsumer(tracer opentracing.Tracer, operationName string, opts ...kitopentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, kitopentracing.WithTags(map[string]interface{}{
		ext.SpanKindConsumer.Key: ext.SpanKindConsumer.Value,
	}))

	return kitopentracing.TraceEndpoint(tracer, operationName, opts...)
}

func TraceProducer(tracer opentracing.Tracer, operationName string, opts ...kitopentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, kitopentracing.WithTags(map[string]interface{}{
		ext.SpanKindProducer.Key: ext.SpanKindProducer.Value,
	}))

	return kitopentracing.TraceEndpoint(tracer, operationName, opts...)
}

type HeadersTextMapCarrier struct {
	headers []kafka.Header
}

func (c HeadersTextMapCarrier) ForeachKey(handler func(key string, val string) error) error {
	for _, h := range c.headers {
		if err := handler(string(h.Key), string(h.Value)); err != nil {
			return err
		}
	}

	return nil
}

func (c *HeadersTextMapCarrier) Set(key, val string) {
	for i := range c.headers {
		if string(c.headers[i].Key) == key {
			c.headers[i] = kafka.Header{
				Key:   []byte(key),
				Value: []byte(val),
			}
			return
		}
	}
	c.headers = append(c.headers, kafka.Header{
		Key:   []byte(key),
		Value: []byte(val),
	})
}

func ContextToKafka(tracer opentracing.Tracer, logger log.Logger) transport.RequestFunc {
	return func(ctx context.Context, msg *kafka.Message) context.Context {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			if err := tracer.Inject(
				span.Context(),
				opentracing.TextMap,
				HeadersTextMapCarrier{
					headers: msg.Headers,
				},
			); err != nil {
				err = fmt.Errorf("failed to inject span context: %w", err)
				_ = logger.Log("err", err)
			}
		}

		return ctx
	}
}

func KafkaToContext(tracer opentracing.Tracer, operationName string, logger log.Logger) transport.RequestFunc {
	return func(ctx context.Context, msg *kafka.Message) context.Context {
		var consumerSpan opentracing.Span
		wireContext, err := tracer.Extract(
			opentracing.TextMap,
			HeadersTextMapCarrier{
				headers: msg.Headers,
			},
		)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			err = fmt.Errorf("failed to extract span context: %w", err)
			_ = logger.Log("err", err)
		}

		consumerSpan = tracer.StartSpan(
			operationName,
			ext.SpanKindConsumer,
			opentracing.ChildOf(wireContext),
		)
		defer consumerSpan.Finish()

		return opentracing.ContextWithSpan(ctx, consumerSpan)
	}
}
