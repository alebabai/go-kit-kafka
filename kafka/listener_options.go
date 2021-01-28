package kafka

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

func ListenerManualCommit(manualCommit bool) ListenerOption {
	return func(l *Listener) {
		l.manualCommit = manualCommit
	}
}
