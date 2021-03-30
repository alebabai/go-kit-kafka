package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/alebabai/go-kit-kafka/kafka"
	kafkatransport "github.com/alebabai/go-kit-kafka/kafka/transport"

	"github.com/alebabai/go-kit-kafka/examples/sarama/consumer"
	"github.com/alebabai/go-kit-kafka/examples/sarama/consumer/adapter"
	"github.com/alebabai/go-kit-kafka/examples/sarama/consumer/endpoint"
	"github.com/alebabai/go-kit-kafka/examples/sarama/consumer/service"
	"github.com/alebabai/go-kit-kafka/examples/sarama/consumer/transport"
	"github.com/alebabai/go-kit-kafka/examples/sarama/domain"
)

func fatal(logger log.Logger, err error) {
	_ = level.Error(logger).Log("err", err)
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
		defer cancel()
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	_ = logger.Log("msg", "initializing services")

	var svc consumer.Service
	{
		storageSvc, err := service.NewStorageService(
			log.With(logger, "component", "storage_service"),
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init storage: %w", err))
		}
		svc = storageSvc
	}

	var endpoints *endpoint.Endpoints
	{
		var err error
		endpoints, err = endpoint.NewEndpoints(svc)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to create endpoints: %w", err))
		}
	}

	_ = logger.Log("msg", "initializing kafka handlers")

	var kafkaHandler kafka.Handler
	{
		var err error
		kafkaHandler, err = transport.NewKafkaHandler(endpoints)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka handler: %w", err))
		}
	}

	_ = logger.Log("msg", "initializing kafka consumer")

	var kafkaConsumerGroup sarama.ConsumerGroup
	{
		cfg := sarama.NewConfig()
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
		cfg.Consumer.Offsets.AutoCommit.Enable = true

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

		kafkaConsumerGroup, err = sarama.NewConsumerGroupFromClient(domain.GroupID, client)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka consumer: %w", err))
		}

		defer func() {
			if err := kafkaConsumerGroup.Close(); err != nil {
				fatal(logger, fmt.Errorf("failed to close kafka consumer group: %w", err))
			}
		}()
	}

	_ = logger.Log("msg", "initializing http handler")

	var httpHandler http.Handler
	{
		var err error
		httpHandler, err = transport.NewHTTPHandler(endpoints)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to create http handler: %w", err))
		}
	}

	errc := make(chan error, 1)

	go func() {
		// use a router in case if there are many topics
		handlers := map[string]kafka.Handler{
			domain.Topic: kafkaHandler,
		}
		router, err := kafkatransport.NewRouter(handlers)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka router: %w", err))
		}

		handler, err := adapter.NewConsumerGroupHandler(router)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka consumer group: %w", err))
		}

		topics := make([]string, 0)
		for topic := range handlers {
			topics = append(topics, topic)
		}

		for {
			errc <- kafkaConsumerGroup.Consume(ctx, topics, handler)
		}
	}()

	go func() {
		errc <- http.ListenAndServe(":8081", httpHandler)
	}()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	_ = logger.Log("msg", "application started")
	_ = logger.Log("msg", "application stopped", "exit", <-errc)
}
