package opentracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type EndpointOptions struct {
	IgnoreBusinessError bool

	GetOperationName func(ctx context.Context, name string) string

	Tags opentracing.Tags

	GetTags func(ctx context.Context) opentracing.Tags
}

type EndpointOption func(*EndpointOptions)

func WithOptions(options EndpointOptions) EndpointOption {
	return func(o *EndpointOptions) {
		*o = options
	}
}

func WithIgnoreBusinessError(ignoreBusinessError bool) EndpointOption {
	return func(o *EndpointOptions) {
		o.IgnoreBusinessError = ignoreBusinessError
	}
}

func WithOperationNameFunc(getOperationName func(ctx context.Context, name string) string) EndpointOption {
	return func(o *EndpointOptions) {
		o.GetOperationName = getOperationName
	}
}

func WithTags(tags opentracing.Tags) EndpointOption {
	return func(o *EndpointOptions) {
		if o.Tags == nil {
			o.Tags = make(opentracing.Tags)
		}

		for key, value := range tags {
			o.Tags[key] = value
		}
	}
}

func WithTagsFunc(getTags func(ctx context.Context) opentracing.Tags) EndpointOption {
	return func(o *EndpointOptions) {
		o.GetTags = getTags
	}
}
