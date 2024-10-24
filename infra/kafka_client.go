package infra

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
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
func (k *kafkaStreamingClient) PublishAsync(ctx context.Context, msgType string, message []byte, blockHeight int64) error {
	topic, err := k.router.GetTopic(msgType)
	if err != nil {
		log.Warn().Any("msgType", msgType).Err(err).Msg("failed to get topic")
		return err
	}

	record := &kgo.Record{
		Topic: topic,
		Value: message,
		Key:   []byte(fmt.Sprintf("%d", blockHeight)), // Use block height as the key to ensure order
	}
	// Async produce
	k.client.Produce(ctx, record, func(record *kgo.Record, err error) {
		if err != nil {
			log.Warn().Err(err).Str("topic", record.Topic).Msg("failed to deliver message")
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

func (k *kafkaStreamingClient) Flush(ctx context.Context) error {
	return k.client.Flush(ctx)
}

func NewFranzClient(seeds []string, user, password string) (*kgo.Client, error) {
	tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
	// Idempotency is enabled by default.
	opts := []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.Dialer(tlsDialer.DialContext),
		kgo.SASL(plain.Auth{
			User: user,
			Pass: password,
		}.AsMechanism()),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()), // Do not change this, as this is the only compression type supported by most kafka clients
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchMaxBytes(16e6),
		kgo.MaxBufferedRecords(1e6),
		kgo.MetadataMaxAge(60 * time.Second),
		kgo.RecordPartitioner(kgo.UniformBytesPartitioner(1e6, false, true, nil)),
	}

	return kgo.NewClient(opts...)
}
