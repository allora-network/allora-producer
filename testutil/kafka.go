package testutil

import (
	"context"
	"errors"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaTestEnv struct {
	ClusterID   string
	Brokers     []string
	Container   *kafka.KafkaContainer
	Producer    *kgo.Client
	Consumer    *kgo.Client
	AdminClient *kadm.Client
}

func NewKafkaTestEnv(ctx context.Context, clusterID string, topics []string) (*KafkaTestEnv, error) {
	container, err := kafka.Run(ctx, "confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID(clusterID),
		testcontainers.WithWaitStrategy(wait.ForLog("Kafka Server started").WithOccurrence(1).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	brokers, err := container.Brokers(ctx)
	if err != nil {
		return nil, err
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchMaxBytes(16e6),
		kgo.MaxBufferedRecords(1e6),
		kgo.MetadataMaxAge(60 * time.Second),
		kgo.RecordPartitioner(kgo.UniformBytesPartitioner(1e6, false, true, nil)),
	}
	producer, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	consumerOpts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topics...),
	}
	consumer, err := kgo.NewClient(consumerOpts...)
	if err != nil {
		return nil, err
	}

	adminClient := kadm.NewClient(producer)

	return &KafkaTestEnv{
		ClusterID:   clusterID,
		Brokers:     brokers,
		Container:   container,
		Producer:    producer,
		Consumer:    consumer,
		AdminClient: adminClient,
	}, nil
}

func (k *KafkaTestEnv) CreateTopics(ctx context.Context, topics ...string) error {
	resp, err := k.AdminClient.CreateTopics(ctx, -1, -1, nil, topics...)
	if err != nil {
		return err
	}
	if len(resp) == 0 {
		return errors.New("no response from create topics")
	}
	return nil
}

func (k *KafkaTestEnv) Close(ctx context.Context) error {
	k.Producer.Close()
	k.Consumer.Close()
	k.AdminClient.Close()
	return k.Container.Terminate(ctx)
}
