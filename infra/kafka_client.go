package infra

import (
	"context"
	"crypto/tls"
	"errors"
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
	client        KafkaClient
	router        domain.TopicRouter
	numPartitions int32
}

// Update the constructor to accept a TopicRouter.
func NewKafkaClient(client KafkaClient, router domain.TopicRouter, numPartitions int32) (domain.StreamingClient, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if router == nil {
		return nil, errors.New("router is nil")
	}

	return &kafkaStreamingClient{
		client:        client,
		router:        router,
		numPartitions: numPartitions,
	}, nil
}

// Publish publishes a message based on its MessageType.
func (k *kafkaStreamingClient) PublishAsync(ctx context.Context, msgType string, message []byte, blockHeight int64) error {
	topic, err := k.router.GetTopic(msgType)
	if err != nil {
		log.Warn().Any("msgType", msgType).Err(err).Msg("failed to get topic")
		return err
	}

	partition := k.getPartition(blockHeight)
	record := &kgo.Record{
		Topic:     topic,
		Value:     message,
		Partition: partition,
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

func (k *kafkaStreamingClient) getPartition(blockHeight int64) int32 {
	return int32(blockHeight % int64(k.numPartitions)) //nolint:gosec // k.numPartitions is a small number, so the modulo will not overflow an int32
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
		kgo.ProducerBatchCompression(kgo.ZstdCompression()),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchMaxBytes(16e6),
		kgo.RecordPartitioner(kgo.ManualPartitioner()),
		kgo.MaxBufferedRecords(1e6),
		kgo.MetadataMaxAge(60 * time.Second),
	}

	return kgo.NewClient(opts...)
}
