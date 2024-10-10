package kafkaclient

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/allora-network/allora-producer/infra"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/kafka"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaClientTestSuite struct {
	suite.Suite
	ctx context.Context

	kafkaContainer *kafka.KafkaContainer
	producer       *kgo.Client
	adminClient    *kadm.Client
	consumer       *kgo.Client
}

const (
	topicName = "test-topic"
)

func (s *KafkaClientTestSuite) SetupSuite() {
	s.ctx = context.Background()

	kafkaContainer, err := kafka.Run(s.ctx,
		"confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID("test-cluster"),
	)
	s.Require().NoError(err)

	s.kafkaContainer = kafkaContainer

	brokers, err := kafkaContainer.Brokers(s.ctx)
	s.Require().NoError(err)

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchMaxBytes(16e6),
		kgo.MaxBufferedRecords(1e6),
		kgo.MetadataMaxAge(60 * time.Second),
		kgo.RecordPartitioner(kgo.UniformBytesPartitioner(1e6, false, true, nil)),
	}
	s.producer, err = kgo.NewClient(opts...)
	s.Require().NoError(err)

	s.adminClient = kadm.NewClient(s.producer)

	consumerOpts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topicName),
	}
	s.consumer, err = kgo.NewClient(consumerOpts...)
	s.Require().NoError(err)
}

func (s *KafkaClientTestSuite) TearDownSuite() {
	s.producer.Close()
	s.consumer.Close()
	err := s.kafkaContainer.Terminate(s.ctx)
	s.Require().NoError(err)
}

func (s *KafkaClientTestSuite) SetupTest() {
	// Add topics to Kafka
	resp, err := s.adminClient.CreateTopics(s.ctx, -1, -1, nil, topicName)
	s.Require().NoError(err)
	s.Require().Len(resp, 1)
}

func (s *KafkaClientTestSuite) TearDownTest() {
	_, err := s.adminClient.DeleteTopics(s.ctx, topicName)
	s.Require().NoError(err)

	err = s.producer.Flush(s.ctx)
	s.Require().NoError(err)

	err = s.consumer.Flush(s.ctx)
	s.Require().NoError(err)
}

func (s *KafkaClientTestSuite) TestPublishAsync() {
	topicRouter := infra.NewTopicRouter(map[string]string{
		"test-type": topicName,
	})

	streamClient, err := infra.NewKafkaClient(s.producer, topicRouter)
	s.Require().NoError(err)

	jsonMsg := json.RawMessage(`{"test": "test"}`)
	err = streamClient.PublishAsync(s.ctx, "test-type", jsonMsg, 1)
	s.Require().NoError(err)
	err = s.producer.Flush(s.ctx)
	s.Require().NoError(err)

	time.Sleep(2 * time.Second)

	// Consume message
	fetches := s.consumer.PollFetches(s.ctx)
	iter := fetches.RecordIter()

	for !iter.Done() {
		record := iter.Next()
		s.Require().Equal(topicName, record.Topic)
		s.Require().Equal(jsonMsg, json.RawMessage(record.Value))
	}
}

func TestKafkaClientTestSuite(t *testing.T) {
	suite.Run(t, new(KafkaClientTestSuite))
}
