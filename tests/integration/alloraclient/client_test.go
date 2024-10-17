package alloraclient

import (
	"context"
	"testing"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/infra"
	"github.com/stretchr/testify/suite"
)

type AlloraClientTestSuite struct {
	suite.Suite

	ctx    context.Context
	client domain.AlloraClientInterface
}

const (
	rpcURL  = "https://allora-rpc.testnet.allora.network"
	timeout = 10 * time.Second
)

func (s *AlloraClientTestSuite) SetupSuite() {
	s.ctx = context.Background()
	client, err := infra.NewAlloraClient(rpcURL, timeout)
	s.Require().NoError(err)
	s.client = client
}

func (s *AlloraClientTestSuite) TestBlockFunctions() {
	height, err := s.client.GetLatestBlockHeight(s.ctx)
	s.Require().NoError(err)
	s.Require().Positive(height)

	// Fetch a previous block to guarantee that the block will be available
	height = height - 10

	header, err := s.client.GetHeader(s.ctx, height)
	s.Require().NoError(err)
	s.Require().NotNil(header)
	s.Require().Equal(height, header.Header.Height)

	block, err := s.client.GetBlockByHeight(s.ctx, height)
	s.Require().NoError(err)
	s.Require().NotNil(block)
	s.Require().Equal(height, block.Block.Header.Height)
	s.Require().Equal(header.Header.Hash(), block.Block.Header.Hash())

	blockResults, err := s.client.GetBlockResults(s.ctx, height)
	s.Require().NoError(err)
	s.Require().NotNil(blockResults)
	s.Require().Equal(height, blockResults.Height)
}

func TestAlloraClientTestSuite(t *testing.T) {
	suite.Run(t, new(AlloraClientTestSuite))
}
