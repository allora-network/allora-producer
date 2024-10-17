package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/app/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewEventsProducer(t *testing.T) {
	service := new(mocks.ProcessorService)
	client := new(mocks.AlloraClientInterface)
	repo := new(mocks.ProcessedBlockRepositoryInterface)

	startHeight := int64(100)
	blockRefreshInterval := time.Millisecond * 50
	rateLimitInterval := time.Millisecond * 50
	numWorkers := 5

	producer, err := NewEventsProducer(service, client, repo, startHeight, blockRefreshInterval, rateLimitInterval, numWorkers)

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

func TestNewEventsProducer_Error(t *testing.T) {
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
		producer, err := NewEventsProducer(tc.service, tc.client, tc.repo, tc.startHeight, tc.blockRefreshInterval, tc.rateLimitInterval, tc.numWorkers)
		require.Error(t, err)
		assert.Nil(t, producer)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestEventsProducer_Execute_InitStartHeightError(t *testing.T) {
	service := new(mocks.ProcessorService)
	client := new(mocks.AlloraClientInterface)
	repo := new(mocks.ProcessedBlockRepositoryInterface)

	producer, err := NewEventsProducer(service, client, repo, 0, time.Millisecond*50, time.Millisecond*50, 5)
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

func TestEventsProducer_ProcessBlock_ErrorFetchingBlock(t *testing.T) {
	service := new(mocks.ProcessorService)
	client := new(mocks.AlloraClientInterface)
	repo := new(mocks.ProcessedBlockRepositoryInterface)

	producer, _ := NewEventsProducer(service, client, repo, 100, time.Millisecond*50, time.Millisecond*50, 5)

	ctx := context.Background()
	height := int64(101)

	client.On("GetBlockResults", ctx, height).Return(nil, errors.New("fetch error"))

	err := producer.processBlockResults(ctx, height)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get block results for height 101: fetch error")

	client.AssertExpectations(t)
}

func TestEventsProducer_ProcessBlock_Success(t *testing.T) {
	service := new(mocks.ProcessorService)
	client := new(mocks.AlloraClientInterface)
	repo := new(mocks.ProcessedBlockRepositoryInterface)

	producer, _ := NewEventsProducer(service, client, repo, 100, time.Millisecond*50, time.Millisecond*50, 5)

	ctx := context.Background()
	height := int64(101)
	expectedResults := &coretypes.ResultBlockResults{
		Height: height,
		TxsResults: []*abci.ExecTxResult{
			{
				Code: 0,
				Data: []byte("result data"),
				Events: []abci.Event{
					{
						Type: "event_type",
						Attributes: []abci.EventAttribute{
							{
								Key:   "key1",
								Value: "value1",
								Index: true,
							},
						},
					},
				},
			},
		},
		FinalizeBlockEvents: []abci.Event{
			{
				Type: "finalize_event_type",
				Attributes: []abci.EventAttribute{
					{
						Key:   "key2",
						Value: "value2",
						Index: true,
					},
				},
			},
		},
	}

	expectedHeader := &coretypes.ResultHeader{
		Header: &types.Header{
			Height: height,
		},
	}

	client.On("GetBlockResults", ctx, height).Return(expectedResults, nil)
	client.On("GetHeader", ctx, height).Return(expectedHeader, nil)
	service.On("ProcessBlockResults", ctx, expectedResults, expectedHeader.Header).Return(nil)
	repo.On("SaveProcessedBlockEvent", ctx, mock.MatchedBy(func(pb domain.ProcessedBlockEvent) bool {
		// Check that Height and Status match expected values
		if pb.Height != height || pb.Status != domain.StatusCompleted {
			return false
		}
		// Ensure ProcessedAt is within the last second
		return time.Since(pb.ProcessedAt) < time.Second
	})).Return(nil)

	err := producer.processBlockResults(ctx, height)
	require.NoError(t, err)

	client.AssertExpectations(t)
	service.AssertExpectations(t)
}
