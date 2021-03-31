package adapter

import (
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
