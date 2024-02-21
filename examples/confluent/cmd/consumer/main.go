package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kafkatransport "github.com/alebabai/go-kit-kafka/v2/transport"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/alebabai/go-kit-kafka/v2/examples/common/consumer"
	"github.com/alebabai/go-kit-kafka/v2/examples/common/domain"

	"github.com/alebabai/go-kit-kafka/v2/examples/confluent/pkg/consumer/kafka/adapter"
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

	_ = logger.Log("msg", "initialization of the application")

	_ = logger.Log("msg", "initialize services")

	var svc consumer.Service
	{
		var err error
		svc, err = consumer.NewStorageService(
			log.With(logger, "component", "storage-service"),
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init storage: %w", err))
		}
	}

	_ = logger.Log("msg", "initialize endpoints")

	var endpoints consumer.Endpoints
	{
		endpoints = consumer.Endpoints{
			CreateEventEndpoint: consumer.MakeCreateEventEndpoint(svc),
			ListEventsEndpoint:  consumer.MakeListEventsEndpoint(svc),
		}
	}

	_ = logger.Log("msg", "initialize kafka handlers")

	kafkaHandler := consumer.NewKafkaHandler(endpoints.CreateEventEndpoint)

	_ = logger.Log("msg", "initialize kafka consumer")

	var kafkaListener *adapter.Listener
	{
		brokerAddr := domain.BrokerAddr
		if v, ok := os.LookupEnv("BROKER_ADDR"); ok {
			brokerAddr = v
		}

		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":  brokerAddr,
			"group.id":           domain.GroupID,
			"enable.auto.commit": true,
		})
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka consumer: %w", err))
		}

		defer func() {
			if err := c.Close(); err != nil {
				fatal(logger, fmt.Errorf("failed to close kafka consumer: %w", err))
			}
		}()

		// use a router in case if there are many topics
		router := make(kafkatransport.Router)
		router.AddHandler(domain.Topic, kafkaHandler)

		topics := make([]string, 0)
		for topic := range router {
			topics = append(topics, topic)
		}

		if err := c.SubscribeTopics(topics, nil); err != nil {
			fatal(logger, fmt.Errorf("failed to subscribe to topics: %w", err))
		}

		kafkaListener, err = adapter.NewListener(
			c,
			router,
			adapter.ListenerErrorHandler(
				transport.NewLogErrorHandler(
					level.Error(
						log.With(logger, "component", "listener"),
					),
				),
			),
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka listener: %w", err))
		}
	}

	_ = logger.Log("msg", "initialize http server")

	var httpServer *http.Server
	{
		httpServer = &http.Server{
			Addr:    ":8081",
			Handler: consumer.NewHTTPHandler(endpoints),
		}

		defer func() {
			if err := httpServer.Shutdown(ctx); err != nil {
				fatal(logger, fmt.Errorf("failed to shutdown http server: %w", err))
			}
		}()
	}

	errc := make(chan error, 1)

	go func() {
		if err := kafkaListener.Listen(ctx); err != nil {
			errc <- err
		}
	}()

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
