package infra

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/util"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

var _ domain.StreamingClient = &kafkaStreamingClient{}

//go:generate mockery --name=KafkaClient
type KafkaClient interface {
	Produce(ctx context.Context, r *kgo.Record, promise func(*kgo.Record, error))
	Flush(ctx context.Context) error
	Close()
}

type kafkaStreamingClient struct {
	client KafkaClient
	router domain.TopicRouter
}

// Update the constructor to accept a TopicRouter.
func NewKafkaClient(client KafkaClient, router domain.TopicRouter) (domain.StreamingClient, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if router == nil {
		return nil, errors.New("router is nil")
	}

	return &kafkaStreamingClient{
		client: client,
		router: router,
	}, nil
}

// Publish publishes a message based on its MessageType.
func (k *kafkaStreamingClient) PublishAsync(ctx context.Context, msgType string, message []byte) error {
	defer util.LogExecutionTime(time.Now(), "Publish")

	topic, err := k.router.GetTopic(msgType)
	if err != nil {
		log.Warn().Any("msgType", msgType).Err(err).Msg("failed to get topic")
		return err
	}

	record := &kgo.Record{
		Topic: topic,
		Value: message,
	}
	// Async produce
	k.client.Produce(ctx, record, func(record *kgo.Record, err error) {
		if err != nil {
			log.Error().Err(err).Str("topic", record.Topic).Msg("failed to deliver message")
		} else {
			log.Debug().Str("topic", record.Topic).Msg("message delivered")
		}
	})
	return nil
}

func (k *kafkaStreamingClient) Close() error {
	// Gracefully flush all pending messages before closing
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	k.client.Flush(ctx)
	k.client.Close()
	return nil
}

func NewFranzClient(seeds []string, user, password string) (*kgo.Client, error) {
	tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
	opts := []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.Dialer(tlsDialer.DialContext),
		kgo.SASL(plain.Auth{
			User: user,
			Pass: password,
		}.AsMechanism()),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.RecordRetries(10),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchMaxBytes(1e6),
		kgo.ProducerLinger(100 * time.Millisecond),
		kgo.RetryTimeout(30 * time.Second),
		kgo.RetryBackoffFn(func(attempt int) time.Duration {
			return time.Duration(attempt) * 100 * time.Millisecond
		}),
	}

	return kgo.NewClient(opts...)
}
