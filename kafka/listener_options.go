package kafka

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
)

type ListenerOption func(*Listener)

func ListenerErrorLogger(logger log.Logger) ListenerOption {
	return func(l *Listener) {
		l.errorHandler = transport.NewLogErrorHandler(logger)
	}
}

func ListenerErrorHandler(errHandler transport.ErrorHandler) ListenerOption {
	return func(l *Listener) {
		l.errorHandler = errHandler
	}
}

func ListenerReadTimeout(readTimeout time.Duration) ListenerOption {
	return func(l *Listener) {
		l.readTimeout = readTimeout
	}
}

func ListenerAsyncHandle(asyncHandle bool) ListenerOption {
	return func(l *Listener) {
		l.asyncHandle = asyncHandle
	}
}
