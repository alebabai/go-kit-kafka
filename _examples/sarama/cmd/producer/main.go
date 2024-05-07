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
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/alebabai/go-kit-kafka/v2/examples/common"
	"github.com/alebabai/go-kit-kafka/v2/examples/common/producer"
	"github.com/alebabai/go-kit-kafka/v2/examples/sarama/pkg/producer/kafka/adapter"
)

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

	_ = logger.Log("msg", "initializing services")

	svc := producer.NewService(
		log.With(logger, "component", "producer-service"),
	)

	_ = logger.Log("msg", "initializing kafka producer")

	var producerMiddleware endpoint.Middleware
	{
		cfg := sarama.NewConfig()
		cfg.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
		cfg.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
		cfg.Producer.Return.Successes = true

		brokerAddr := common.BrokerAddr
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
			common.KafkaTopic,
		).Endpoint()

		producerMiddleware = producer.Middleware(e)
	}

	_ = logger.Log("msg", "initializing http server")

	var httpServer *http.Server
	{
		httpServer = &http.Server{
			Addr: producer.HTTPServerAddr,
			Handler: producer.NewHTTPHandler(
				producerMiddleware(
					producer.MakeGenerateEventEndpoint(svc),
				),
			),
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

func fatal(logger log.Logger, err error) {
	_ = logger.Log("err", err)
	os.Exit(1)
}
