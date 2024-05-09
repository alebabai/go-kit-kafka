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
	"github.com/alebabai/go-kafka"
	adapter "github.com/alebabai/go-kafka/adapter/sarama"
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

	_ = logger.Log("msg", "initializing kafka consumer group")

	var consumerGroupListener *adapter.ConsumerGroupListener
	{
		cfg := sarama.NewConfig()
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
		cfg.Consumer.Offsets.AutoCommit.Enable = true
		cfg.Consumer.Return.Errors = true

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

		cg, err := sarama.NewConsumerGroupFromClient(consumer.KafkaGroupID, client)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka consumer group: %w", err))
		}

		defer func() {
			if err := cg.Close(); err != nil {
				fatal(logger, fmt.Errorf("failed to close kafka consumer group: %w", err))
			}
		}()

		// use a router in case if there are many topics
		router := kafka.StaticRouterByTopic{
			common.KafkaTopic: consumer.NewKafkaHandler(
				consumer.MakeCreateEventEndpoint(svc),
			),
		}

		cgh, err := adapter.NewConsumerGroupHandler(router)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka consumer group handler: %w", err))
		}

		consumerGroupListener, err = adapter.NewConsumerGroupListener(
			cg,
			cgh,
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init kafka listener: %w", err))
		}
	}

	_ = logger.Log("msg", "initializing http server")

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
		if err := consumerGroupListener.Listen(ctx, common.KafkaTopic); err != nil {
			errc <- err
		}
	}()

	go func() {
		for err := range consumerGroupListener.Errors() {
			time.Sleep(2 * time.Second) // debounce
			level.Error(logger).Log("err", err)
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

func fatal(logger log.Logger, err error) {
	_ = level.Error(logger).Log("err", err)
	os.Exit(1)
}
