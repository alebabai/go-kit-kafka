package transport

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
)

type ConsumerOption func(consumer *Consumer)

func ConsumerBefore(before ...ConsumerRequestFunc) ConsumerOption {
	return func(c *Consumer) {
		c.before = append(c.before, before...)
	}
}

func ConsumerAfter(after ...ConsumerResponseFunc) ConsumerOption {
	return func(c *Consumer) {
		c.after = append(c.after, after...)
	}
}

func ConsumerErrorLogger(logger log.Logger) ConsumerOption {
	return func(c *Consumer) {
		c.errorHandler = transport.NewLogErrorHandler(logger)
	}
}

func ConsumerErrorHandler(errorHandler transport.ErrorHandler) ConsumerOption {
	return func(c *Consumer) {
		c.errorHandler = errorHandler
	}
}
