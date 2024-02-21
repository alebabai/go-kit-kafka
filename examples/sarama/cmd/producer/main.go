package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/alebabai/go-kit-kafka/v2/examples/common/domain"
	"github.com/alebabai/go-kit-kafka/v2/examples/common/producer"
	"github.com/alebabai/go-kit-kafka/v2/examples/sarama/pkg/producer/kafka/adapter"
)

func fatal(logger log.Logger, err error) {
	_ = logger.Log("err", err)
	os.Exit(1)
}

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	{
		ctx = context.Background()
		ctx, cancel = context.WithCancel(ctx)
		defer func() {
			cancel()
			time.Sleep(3 * time.Second)
		}()
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	_ = logger.Log("msg", "initialization of the application")

	_ = logger.Log("msg", "initialize services")

	var svc producer.Service
	{
		var err error
		svc, err = producer.NewGeneratorService(logger)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to create generator: %w", err))
		}
	}

	_ = logger.Log("msg", "initialize kafka producer")

	var producerMiddleware kitendpoint.Middleware
	{
		cfg := sarama.NewConfig()
		cfg.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
		cfg.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
		cfg.Producer.Return.Successes = true

		brokerAddr := domain.BrokerAddr
		if v, ok := os.LookupEnv("BROKER_ADDR"); ok {
			brokerAddr = v
		}

		client, err := sarama.NewClient(
			[]string{brokerAddr},
			cfg,
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka client: %w", err))
		}

		p, err := sarama.NewSyncProducerFromClient(client)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka producer: %w", err))
		}

		defer func() {
			if err := p.Close(); err != nil {
				fatal(logger, fmt.Errorf("failed to close kafka producer: %w", err))
			}
		}()

		e := producer.NewKafkaProducer(
			adapter.NewProducer(p),
			domain.Topic,
		).Endpoint()

		producerMiddleware = producer.Middleware(e)
	}

	_ = logger.Log("msg", "initialize endpoints")

	var endpoints producer.Endpoints
	{
		endpoints = producer.Endpoints{
			GenerateEvent: producerMiddleware(
				producer.MakeGenerateEventEndpoint(svc),
			),
		}
	}

	_ = logger.Log("msg", "initialize http server")

	var httpServer *http.Server
	{
		httpServer = &http.Server{
			Addr:    ":8080",
			Handler: producer.NewHTTPHandler(endpoints),
		}

		defer func() {
			if err := httpServer.Shutdown(ctx); err != nil {
				fatal(logger, fmt.Errorf("failed to shutdown http server: %w", err))
			}
		}()
	}

	errc := make(chan error, 1)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errc <- err
		}
	}()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	_ = logger.Log("msg", "application started")
	_ = logger.Log("msg", "application stopped", "exit", <-errc)
}
