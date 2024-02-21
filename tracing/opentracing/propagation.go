package opentracing

import (
	"context"
	"fmt"

	"github.com/alebabai/go-kafka"
	"github.com/alebabai/go-kit-kafka/v2/transport"
	"github.com/go-kit/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// KafkaHeadersCarrier satisfies both [opentracing.TextMapWriter] and [opentracing.TextMapReader].
//
// Example usage for consumer side:
//
//	carrier := KafkaHeadersCarrier(msg.Headers)
//	clientContext, err := tracer.Extract(opentracing.TextMap, carrier)
//
// Example usage for producer side:
//
//	carrier := KafkaHeadersCarrier(msg.Headers)
//	err := tracer.Inject(
//	    span.Context(),
//	    opentracing.TextMap,
//	    carrier,
//	)
type KafkaHeadersCarrier []kafka.Header

// Set conforms to the opentracing.TextMapWriter interface.
func (c *KafkaHeadersCarrier) Set(key, val string) {
	for i := range *c {
		if string((*c)[i].Key) == key {
			(*c)[i] = kafka.Header{
				Key:   []byte(key),
				Value: []byte(val),
			}
			return
		}
	}
	*c = append(*c, kafka.Header{
		Key:   []byte(key),
		Value: []byte(val),
	})
}

// ForeachKey conforms to the [opentracing.TextMapReader] interface.
func (c KafkaHeadersCarrier) ForeachKey(handler func(key string, val string) error) error {
	for _, h := range c {
		if err := handler(string(h.Key), string(h.Value)); err != nil {
			return err
		}
	}

	return nil
}

// ContextToKafka returns an [transport.RequestFunc] that injects an [opentracing.Span]
// found in ctx into the kafka headers. If no such Span can be found, the
// transport.RequestFunc is a noop.
func ContextToKafka(tracer opentracing.Tracer, logger log.Logger) transport.RequestFunc {
	return func(ctx context.Context, msg *kafka.Message) context.Context {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			if err := tracer.Inject(
				span.Context(),
				opentracing.TextMap,
				KafkaHeadersCarrier(msg.Headers),
			); err != nil {
				err = fmt.Errorf("failed to inject span context: %w", err)
				_ = logger.Log("err", err)
			}
		}

		return ctx
	}
}

// KafkaToContext returns an [transport.RequestFunc] that tries to join with an
// OpenTracing trace found in msg and starts a new Span called
// operationName accordingly. If no trace could be found in msg, the Span
// will be a trace root. The Span is incorporated in the returned [context.Context] and
// can be retrieved with opentracing.SpanFromContext(ctx).
func KafkaToContext(tracer opentracing.Tracer, operationName string, logger log.Logger) transport.RequestFunc {
	return func(ctx context.Context, msg *kafka.Message) context.Context {
		var consumerSpan opentracing.Span
		wireContext, err := tracer.Extract(
			opentracing.TextMap,
			KafkaHeadersCarrier(msg.Headers),
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
