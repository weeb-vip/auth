package handlers

import (
	"context"
	"fmt"
	"strings"

	epKafka "github.com/ThatCatDev/ep/v2/drivers/kafka"
	epNats "github.com/ThatCatDev/ep/v2/drivers/nats"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/weeb-vip/auth/config"
	logger2 "github.com/weeb-vip/auth/internal/logger"
)

// Publish sends one encoded event to the configured destination.
//
// It takes the value rather than a driver message because the two transports
// name that field differently -- Value on kafka.Message, Data on the NATS one --
// and every call site here only ever set the value anyway.
type Publish func(ctx context.Context, value []byte) error

// newPublisher builds the producer for whichever transport EventTransport
// names, returning the publish function and a close for it.
//
// A runtime switch rather than a separate command, because this producer lives
// inside the GraphQL handler: there is no command name to pick a transport with
// the way the standalone consumers have. Staging can be flipped by setting
// PRODUCER_TYPE=nats on the deployment while production stays on the default.
//
// Unknown values are rejected rather than silently falling back to Kafka. A
// typo in the environment would otherwise look exactly like a working service
// that publishes to the wrong place.
func newPublisher(conf config.Config) (Publish, func(), error) {
	switch strings.ToLower(strings.TrimSpace(conf.ProducerType)) {
	case "", "kafka":
		driver := epKafka.NewKafkaDriver(&epKafka.KafkaConfig{
			ConsumerGroupName:       conf.KafkaConfig.ConsumerGroupName,
			BootstrapServers:        conf.KafkaConfig.BootstrapServers,
			ConsumerAutoOffsetReset: &conf.KafkaConfig.Offset,
		})

		topic := conf.KafkaConfig.ProducerTopic

		return func(ctx context.Context, value []byte) error {
				log := logger2.FromCtx(ctx)
				log.Info().Str("topic", topic).Msg("Producing message to Kafka")
				if err := driver.Produce(ctx, topic, &kafka.Message{Value: value}); err != nil {
					log.Error().Err(err).Str("topic", topic).Msg("Failed to produce message")

					return err
				}

				return nil
			}, func() {
				_ = driver.Close()
			}, nil

	case "nats":
		driver := epNats.NewNatsDriver(&epNats.Config{
			URL:        conf.NatsConfig.URL,
			StreamName: conf.NatsConfig.StreamName,
		})

		subject := conf.NatsConfig.ProducerSubject

		return func(ctx context.Context, value []byte) error {
				log := logger2.FromCtx(ctx)
				log.Info().Str("subject", subject).Msg("Producing message to NATS")
				if err := driver.Produce(ctx, subject, &epNats.Message{Data: value}); err != nil {
					log.Error().Err(err).Str("subject", subject).Msg("Failed to produce message")

					return err
				}

				return nil
			}, func() {
				_ = driver.Close()
			}, nil

	default:
		return nil, nil, fmt.Errorf("unknown PRODUCER_TYPE %q: expected kafka or nats", conf.ProducerType)
	}
}
