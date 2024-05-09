package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alebabai/go-kafka"
	adapter "github.com/alebabai/go-kafka/adapter/confluent"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/alebabai/go-kit-kafka/v2/examples/common"
	"github.com/alebabai/go-kit-kafka/v2/examples/common/consumer"
)

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

	_ = logger.Log("msg", "initializing services")

	svc := consumer.NewService(
		log.With(logger, "component", "consumer-service"),
	)

	_ = logger.Log("msg", "initialize kafka consumer")

	var consumerListener *adapter.ConsumerListener
	{
		brokerAddr := common.BrokerAddr
		if v, ok := os.LookupEnv("BROKER_ADDR"); ok {
			brokerAddr = v
		}

		c, err := ckafka.NewConsumer(&ckafka.ConfigMap{
			"bootstrap.servers":  brokerAddr,
			"group.id":           consumer.KafkaGroupID,
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
		router := kafka.StaticRouterByTopic{
			common.KafkaTopic: consumer.NewKafkaHandler(
				consumer.MakeCreateEventEndpoint(svc),
			),
		}
		if err := c.SubscribeTopics([]string{common.KafkaTopic}, nil); err != nil {
			fatal(logger, fmt.Errorf("failed to subscribe to topics: %w", err))
		}

		consumerListener, err = adapter.NewConsumerListener(
			c,
			router,
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka listener: %w", err))
		}
	}

	_ = logger.Log("msg", "initialize http server")

	var httpServer *http.Server
	{
		httpServer = &http.Server{
			Addr:    consumer.HTTPServerAddr,
			Handler: consumer.NewHTTPHandler(consumer.MakeListEventsEndpoint(svc)),
		}

		defer func() {
			if err := httpServer.Shutdown(ctx); err != nil {
				fatal(logger, fmt.Errorf("failed to shutdown http server: %w", err))
			}
		}()
	}

	errc := make(chan error, 1)

	go func() {
		if err := consumerListener.Listen(ctx, -1); err != nil {
			errc <- err
		}
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errc <- err
		}
	}()

	sigc := make(chan os.Signal, 1)

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	}()

	_ = logger.Log("msg", "application started")

	select {
	case sig := <-sigc:
		_ = logger.Log("msg", "application stopped", "exit", sig)
	case err := <-errc:
		fatal(logger, err)
	}
}

func fatal(logger log.Logger, err error) {
	_ = level.Error(logger).Log("msg", "application stopped by an error", "err", err)
	os.Exit(1)
}
