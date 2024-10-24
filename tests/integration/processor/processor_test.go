package processor

import (
	"context"
	"testing"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/app/service"
	"github.com/allora-network/allora-producer/codec"
	"github.com/allora-network/allora-producer/infra"
	"github.com/allora-network/allora-producer/testutil"
	"github.com/stretchr/testify/suite"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
)

const (
	topicName = "test-topic"
	clusterID = "test-cluster"
)

type ProcessorTestSuite struct {
	suite.Suite
	ctx          context.Context
	kafkaTestEnv *testutil.KafkaTestEnv

	topicRouter domain.TopicRouter
	kafkaClient domain.StreamingClient
	codec       domain.CodecInterface

	filterEvent              *service.FilterEvent
	filterTransactionMessage *service.FilterTransactionMessage

	processorService *service.ProcessorService
}

func (s *ProcessorTestSuite) SetupSuite() {
	s.ctx = context.Background()

	kafkaTestEnv, err := testutil.NewKafkaTestEnv(s.ctx, clusterID, []string{topicName})
	s.Require().NoError(err)
	s.kafkaTestEnv = kafkaTestEnv

	topicRouter := infra.NewTopicRouter(map[string]string{
		"test-type": topicName,
	})
	s.topicRouter = topicRouter

	kafkaClient, err := infra.NewKafkaClient(kafkaTestEnv.Producer, topicRouter)
	s.Require().NoError(err)
	s.kafkaClient = kafkaClient

	s.codec = codec.NewCodec()

	s.filterEvent = service.NewFilterEvent("test-type")
	s.filterTransactionMessage = service.NewFilterTransactionMessage("test-type")

	processorService, err := service.NewProcessorService(s.kafkaClient, s.codec, s.filterEvent, s.filterTransactionMessage)
	s.Require().NoError(err)
	s.processorService = processorService
}

func (s *ProcessorTestSuite) TearDownSuite() {
	s.kafkaTestEnv.Close(s.ctx)
}

func (s *ProcessorTestSuite) SetupTest() {
	err := s.kafkaTestEnv.CreateTopics(s.ctx, topicName)
	s.Require().NoError(err)
}

func (s *ProcessorTestSuite) TestProcessBlock() {
	block := &ctypes.ResultBlock{
		BlockID: types.BlockID{
			Hash: []byte("test-hash"),
			PartSetHeader: types.PartSetHeader{
				Total: 1,
				Hash:  []byte("test-hash"),
			},
		},
		Block: &types.Block{
			Header: types.Header{
				Height: 1,
			},
			Data: types.Data{
				Txs: []types.Tx{
					[]byte("test-tx"),
				},
			},
		},
	}

	err := s.processorService.ProcessBlock(s.ctx, block)
	s.Require().NoError(err)
}

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}
