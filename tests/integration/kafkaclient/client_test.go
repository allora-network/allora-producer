package kafkaclient

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/allora-network/allora-producer/infra"
	"github.com/allora-network/allora-producer/testutil"
	"github.com/stretchr/testify/suite"
)

type KafkaClientTestSuite struct {
	suite.Suite
	ctx context.Context

	kafkaTestEnv *testutil.KafkaTestEnv
}

const (
	topicName = "test-topic"
	clusterID = "test-cluster"
)

func (s *KafkaClientTestSuite) SetupSuite() {
	s.ctx = context.Background()

	kafkaTestEnv, err := testutil.NewKafkaTestEnv(s.ctx, clusterID, []string{topicName})
	s.Require().NoError(err)

	s.kafkaTestEnv = kafkaTestEnv
}

func (s *KafkaClientTestSuite) TearDownSuite() {
	s.kafkaTestEnv.Close(s.ctx)
}

func (s *KafkaClientTestSuite) SetupTest() {
	// Add topics to Kafka
	err := s.kafkaTestEnv.CreateTopics(s.ctx, topicName)
	s.Require().NoError(err)
}

func (s *KafkaClientTestSuite) TearDownTest() {
	err := s.kafkaTestEnv.Producer.Flush(s.ctx)
	s.Require().NoError(err)

	err = s.kafkaTestEnv.Consumer.Flush(s.ctx)
	s.Require().NoError(err)

	resp, err := s.kafkaTestEnv.AdminClient.DeleteTopics(s.ctx, topicName)
	s.Require().NoError(err)
	s.Require().Len(resp, 1)
}

func (s *KafkaClientTestSuite) TestPublishAsync() {
	topicRouter := infra.NewTopicRouter(map[string]string{
		"test-type": topicName,
	})

	client, err := infra.NewKafkaClient(s.kafkaTestEnv.Producer, topicRouter)
	s.Require().NoError(err)

	jsonMsg := json.RawMessage(`{"test": "test"}`)
	err = client.PublishAsync(s.ctx, "test-type", jsonMsg, 1)
	s.Require().NoError(err)
	err = s.kafkaTestEnv.Producer.Flush(s.ctx)
	s.Require().NoError(err)

	time.Sleep(2 * time.Second)

	// Consume message
	fetches := s.kafkaTestEnv.Consumer.PollFetches(s.ctx)
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
