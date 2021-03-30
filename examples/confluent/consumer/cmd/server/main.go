package main

import (
	"context"
	"fmt"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alebabai/go-kit-kafka/kafka"

	"github.com/alebabai/go-kit-kafka/examples/confluent/consumer"
	"github.com/alebabai/go-kit-kafka/examples/confluent/consumer/adapter"
	"github.com/alebabai/go-kit-kafka/examples/confluent/consumer/endpoint"
	"github.com/alebabai/go-kit-kafka/examples/confluent/consumer/service"
	"github.com/alebabai/go-kit-kafka/examples/confluent/consumer/transport"
	"github.com/alebabai/go-kit-kafka/examples/confluent/domain"
)

func fatal(logger log.Logger, err error) {
	_ = level.Error(logger).Log("err: %w", err)
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

	var httpHandler http.Handler
	{
		var err error
		httpHandler, err = transport.NewHTTPHandler(endpoints)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to create http handler: %w", err))
		}
	}

	var kafkaHandler kafka.Handler
	{
		var err error
		kafkaHandler, err = transport.NewKafkaHandler(endpoints)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka handler: %w", err))
		}
	}

	var kafkaListener *kafka.Listener
	{
		brokerAddr := domain.BrokerAddr
		if v, ok := os.LookupEnv("BROKER_ADDR"); ok {
			brokerAddr = v
		}

		c, err := ckafka.NewConsumer(&ckafka.ConfigMap{
			"bootstrap.servers":  brokerAddr,
			"group.id":           domain.GroupID,
			"enable.auto.commit": true,
			//"go.events.channel.enable":        true,
			//"go.application.rebalance.enable": true,
		})
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka c: %w", err))
		}

		topics := []string{
			domain.Topic,
		}
		if err := c.SubscribeTopics(topics, nil); err != nil {
			fatal(logger, fmt.Errorf("failed to subscribe to topics: %w", err))
		}

		reader, err := adapter.NewFunctionReader(c)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka reader: %w", err))
		}

		handlers := map[string]kafka.Handler{
			domain.Topic: kafkaHandler,
		}

		kafkaListener, err = kafka.NewListener(
			reader,
			handlers,
			kafka.ListenerErrorLogger(
				log.With(logger, "component", "listener"),
			),
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka listener: %w", err))
		}
	}

	errc := make(chan error, 1)

	go func() {
		errc <- http.ListenAndServe(":8081", httpHandler)
	}()

	go func() {
		errc <- kafkaListener.Listen(ctx)
	}()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	_ = logger.Log("exit", <-errc)
}
