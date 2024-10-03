package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/app/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewTransactionsProducer(t *testing.T) {
	service := new(mocks.ProcessorService)
	client := new(mocks.AlloraClientInterface)
	repo := new(mocks.ProcessedBlockRepositoryInterface)

	startHeight := int64(100)
	blockRefreshInterval := time.Millisecond * 50
	rateLimitInterval := time.Millisecond * 50
	numWorkers := 5

	producer, err := NewTransactionsProducer(service, client, repo, startHeight, blockRefreshInterval, rateLimitInterval, numWorkers)

	require.NoError(t, err)
	assert.NotNil(t, producer)
	assert.Equal(t, service, producer.service)
	assert.Equal(t, client, producer.alloraClient)
	assert.Equal(t, repo, producer.repository)
	assert.Equal(t, startHeight, producer.startHeight)
	assert.Equal(t, blockRefreshInterval, producer.blockRefreshInterval)
	assert.Equal(t, rateLimitInterval, producer.rateLimitInterval)
	assert.Equal(t, numWorkers, producer.numWorkers)
}

func TestNewTransactionsProducer_Error(t *testing.T) {
	testCases := []struct {
		service              domain.ProcessorService
		client               domain.AlloraClientInterface
		repo                 domain.ProcessedBlockRepositoryInterface
		startHeight          int64
		blockRefreshInterval time.Duration
		rateLimitInterval    time.Duration
		numWorkers           int
		expectedError        error
	}{
		{nil, new(mocks.AlloraClientInterface), new(mocks.ProcessedBlockRepositoryInterface), 0, 0, 0, 0, errors.New("service is nil")},
		{new(mocks.ProcessorService), nil, new(mocks.ProcessedBlockRepositoryInterface), 0, 0, 0, 0, errors.New("client is nil")},
		{new(mocks.ProcessorService), new(mocks.AlloraClientInterface), nil, 0, 0, 0, 0, errors.New("repository is nil")},
	}

	for _, tc := range testCases {
		producer, err := NewTransactionsProducer(tc.service, tc.client, tc.repo, tc.startHeight, tc.blockRefreshInterval, tc.rateLimitInterval, tc.numWorkers)
		require.Error(t, err)
		assert.Nil(t, producer)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestTransactionsProducer_Execute_InitStartHeightError(t *testing.T) {
	service := new(mocks.ProcessorService)
	client := new(mocks.AlloraClientInterface)
	repo := new(mocks.ProcessedBlockRepositoryInterface)

	producer, err := NewTransactionsProducer(service, client, repo, 0, time.Millisecond*50, time.Millisecond*50, 5)
	require.NoError(t, err)

	// Simulate InitStartHeight returning an error
	client.On("GetLatestBlockHeight", mock.Anything).Return(int64(0), errors.New("init error")).Times(6)

	err = producer.Execute(context.Background())
	require.Error(t, err)
	assert.Equal(t, "failed to get latest block height after 5 retries: init error", err.Error())
}

// func TestTransactionsProducer_Execute_Success(t *testing.T) {
// 	service := new(mocks.ProcessorService)
// 	client := new(mocks.AlloraClientInterface)
// 	repo := new(mocks.ProcessedBlockRepositoryInterface)

// 	producer, _ := NewTransactionsProducer(service, client, repo, 100, time.Second*10, time.Second, 5)

// 	// Mock InitStartHeight to succeed
// 	producer.BaseProducer.InitStartHeight = func(ctx context.Context) error {
// 		return nil
// 	}

// 	// Mock MonitorLoopParallel to succeed
// 	producer.BaseProducer.MonitorLoopParallel = func(ctx context.Context, processBlockFunc func(context.Context, int64) error, workers int) error {
// 		return nil
// 	}

// 	err := producer.Execute(context.Background())
// 	assert.NoError(t, err)
// }

func TestTransactionsProducer_ProcessBlock_ErrorFetchingBlock(t *testing.T) {
	service := new(mocks.ProcessorService)
	client := new(mocks.AlloraClientInterface)
	repo := new(mocks.ProcessedBlockRepositoryInterface)

	producer, _ := NewTransactionsProducer(service, client, repo, 100, time.Millisecond*50, time.Millisecond*50, 5)

	ctx := context.Background()
	height := int64(101)

	client.On("GetBlockByHeight", ctx, height).Return(nil, errors.New("fetch error"))

	err := producer.processBlock(ctx, height)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get block for height 101: fetch error")

	client.AssertExpectations(t)
}

func TestTransactionsProducer_ProcessBlock_Success(t *testing.T) {
	service := new(mocks.ProcessorService)
	client := new(mocks.AlloraClientInterface)
	repo := new(mocks.ProcessedBlockRepositoryInterface)

	producer, _ := NewTransactionsProducer(service, client, repo, 100, time.Millisecond*50, time.Millisecond*50, 5)

	ctx := context.Background()
	height := int64(101)
	block := &ctypes.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height: height,
			},
			Data: types.Data{
				Txs: types.Txs{
					[]byte("tx1"),
					[]byte("tx2"),
				},
			},
		},
	}

	client.On("GetBlockByHeight", ctx, height).Return(block, nil)
	service.On("ProcessBlock", ctx, block).Return(nil)
	repo.On("SaveProcessedBlock", ctx, mock.MatchedBy(func(pb domain.ProcessedBlock) bool {
		// Check that Height and Status match expected values
		if pb.Height != height || pb.Status != domain.StatusCompleted {
			return false
		}
		// Ensure ProcessedAt is within the last second
		return time.Since(pb.ProcessedAt) < time.Second
	})).Return(nil)

	err := producer.processBlock(ctx, height)
	require.NoError(t, err)

	client.AssertExpectations(t)
	service.AssertExpectations(t)
}
