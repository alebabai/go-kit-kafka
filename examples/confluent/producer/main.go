package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
)

func fatal(logger log.Logger, err error) {
	_ = logger.Log("err", err)
	os.Exit(1)
}

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	var p *kafka.Producer
	{
		brokerAddr := domain.BrokerAddr
		if v, ok := os.LookupEnv("BROKER_ADDR"); ok {
			brokerAddr = v
		}

		var err error
		p, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": brokerAddr,
		})
		if err != nil {
			fatal(logger, fmt.Errorf("failed to create publisher: %w", err))
		}

		defer p.Close()
	}

	var svc Service
	{
		g, err := NewGenerator(logger)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to create generator: %w", err))
		}

		svc = ProducerMiddleware(domain.Topic, p, logger)(g)
	}

	var e *Endpoints
	{
		var err error
		e, err = NewEndpoints(svc)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to create endpoints: %w", err))
		}
	}

	var h http.Handler
	{
		var err error
		h, err = NewHTTPHandler(e)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to create http handler: %w", err))
		}
	}

	errc := make(chan error, 1)

	go func() {
		errc <- http.ListenAndServe(":8080", h)
	}()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	_ = logger.Log("exit", <-errc)
}
