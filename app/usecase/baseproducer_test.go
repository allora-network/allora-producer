package usecase

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/allora-network/allora-producer/app/domain/mocks"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBaseProducer_InitStartHeight(t *testing.T) {
	mockAlloraClient := new(mocks.AlloraClientInterface)
	mockRepository := new(mocks.ProcessedBlockRepositoryInterface)
	mockService := new(mocks.ProcessorService)
	logger := zerolog.Nop()

	bm := &BaseProducer{
		alloraClient:      mockAlloraClient,
		repository:        mockRepository,
		service:           mockService,
		logger:            &logger,
		rateLimitInterval: 10 * time.Millisecond,
	}

	t.Run("StartHeight is zero, fetch latest height", func(t *testing.T) {
		mockAlloraClient.On("GetLatestBlockHeight", mock.Anything).Return(int64(100), nil)

		bm.startHeight = 0
		err := bm.InitStartHeight(context.Background())

		require.NoError(t, err)
		require.Equal(t, int64(100), bm.startHeight)
		mockAlloraClient.AssertCalled(t, "GetLatestBlockHeight", mock.Anything)
	})

	// Reset mockAlloraClient
	t.Run("StartHeight is non-zero, do not fetch latest height", func(t *testing.T) {
		bm.startHeight = 50

		mockAlloraClient.ExpectedCalls = nil
		mockAlloraClient.Calls = nil

		err := bm.InitStartHeight(context.Background())

		require.NoError(t, err)
		require.Equal(t, int64(50), bm.startHeight)
		mockAlloraClient.AssertNotCalled(t, "GetLatestBlockHeight", mock.Anything)
	})

	t.Run("Error fetching latest height", func(t *testing.T) {
		mockAlloraClient.ExpectedCalls = nil
		mockAlloraClient.Calls = nil

		mockAlloraClient.On("GetLatestBlockHeight", mock.Anything).Return(int64(0), errors.New("network error"))

		bm.startHeight = 0
		err := bm.InitStartHeight(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get latest block height")
		mockAlloraClient.AssertCalled(t, "GetLatestBlockHeight", mock.Anything)
	})

	t.Run("StartHeight is zero, fetch latest height after 5 retries", func(t *testing.T) {
		bm.startHeight = 0
		mockAlloraClient.ExpectedCalls = nil
		mockAlloraClient.Calls = nil

		mockAlloraClient.On("GetLatestBlockHeight", mock.Anything).Return(int64(0), errors.New("network error")).Times(5)
		mockAlloraClient.On("GetLatestBlockHeight", mock.Anything).Return(int64(100), nil)

		bm.startHeight = 0
		err := bm.InitStartHeight(context.Background())

		require.NoError(t, err)
		require.Equal(t, int64(100), bm.startHeight)
		mockAlloraClient.AssertCalled(t, "GetLatestBlockHeight", mock.Anything)
	})
}

func TestBaseProducer_MonitorLoop(t *testing.T) {
	mockAlloraClient := new(mocks.AlloraClientInterface)
	mockRepository := new(mocks.ProcessedBlockRepositoryInterface)
	mockService := new(mocks.ProcessorService)
	logger := zerolog.Nop()

	bm := &BaseProducer{
		alloraClient:         mockAlloraClient,
		repository:           mockRepository,
		service:              mockService,
		logger:               &logger,
		startHeight:          1,
		blockRefreshInterval: 10 * time.Millisecond,
		rateLimitInterval:    10 * time.Millisecond,
	}

	t.Run("Process blocks successfully", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		processBlock := func(_ context.Context, height int64) error { //nolint:unparam
			require.Equal(t, bm.startHeight, height)
			return nil
		}

		var wg sync.WaitGroup
		wg.Add(1)
		// Run MonitorLoop in a separate goroutine
		go func() {
			err := bm.MonitorLoop(ctx, processBlock)
			assert.EqualError(t, err, "context canceled")
			wg.Done()
		}()

		// Let it run for a short time
		time.Sleep(50 * time.Millisecond)
		cancel()

		require.GreaterOrEqual(t, bm.startHeight, int64(3))
		wg.Wait()
	})

	t.Run("Process block with error and retry", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		callCount := 0
		processBlock := func(_ context.Context, _ int64) error {
			if callCount < 1 {
				callCount++
				return errors.New("processing error")
			}
			return nil
		}

		var wg sync.WaitGroup
		wg.Add(1)

		// Run MonitorLoop in a separate goroutine
		go func() {
			err := bm.MonitorLoop(ctx, processBlock)
			assert.EqualError(t, err, "context canceled")
			wg.Done()
		}()

		// Let it run for a short time
		time.Sleep(50 * time.Millisecond)
		cancel()

		require.GreaterOrEqual(t, bm.startHeight, int64(3))
		require.Equal(t, 1, callCount)
		wg.Wait()
	})
}

func TestBaseProducer_MonitorLoopParallel(t *testing.T) {
	mockAlloraClient := new(mocks.AlloraClientInterface)
	mockRepository := new(mocks.ProcessedBlockRepositoryInterface)
	mockService := new(mocks.ProcessorService)
	logger := zerolog.Nop()

	bm := &BaseProducer{
		alloraClient:         mockAlloraClient,
		repository:           mockRepository,
		service:              mockService,
		logger:               &logger,
		startHeight:          1,
		blockRefreshInterval: 10 * time.Millisecond,
		rateLimitInterval:    10 * time.Millisecond,
	}

	mockAlloraClient.On("GetLatestBlockHeight", mock.Anything).Return(int64(5), nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mu sync.Mutex
	processed := []int64{}

	processBlock := func(_ context.Context, height int64) error { //nolint:unparam
		mu.Lock()
		processed = append(processed, height)
		mu.Unlock()
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := bm.MonitorLoopParallel(ctx, processBlock, 2)
		assert.NoError(t, err)
		wg.Done()
	}()

	// Let it run for a short time
	time.Sleep(100 * time.Millisecond)
	cancel()

	require.ElementsMatch(t, []int64{1, 2, 3, 4, 5}, processed)
	mockAlloraClient.AssertCalled(t, "GetLatestBlockHeight", mock.Anything)
	wg.Wait()
}

func TestBaseProducer_MonitorLoopParallel_ErrorProcessBlock(t *testing.T) {
	mockAlloraClient := new(mocks.AlloraClientInterface)
	mockRepository := new(mocks.ProcessedBlockRepositoryInterface)
	mockService := new(mocks.ProcessorService)
	logger := zerolog.Nop()

	bm := &BaseProducer{
		alloraClient:         mockAlloraClient,
		repository:           mockRepository,
		service:              mockService,
		logger:               &logger,
		startHeight:          1,
		blockRefreshInterval: 10 * time.Millisecond,
		rateLimitInterval:    10 * time.Millisecond,
	}

	mockAlloraClient.On("GetLatestBlockHeight", mock.Anything).Return(int64(2), nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mu sync.Mutex
	processed := []int64{}
	errorsCount := make(map[int64]int)

	processBlock := func(_ context.Context, height int64) error {
		mu.Lock()
		defer mu.Unlock()
		processed = append(processed, height)
		errorsCount[height]++
		if errorsCount[height] < 2 {
			return errors.New("processing error")
		}
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := bm.MonitorLoopParallel(ctx, processBlock, 2)
		assert.NoError(t, err)
		wg.Done()
	}()

	// Let it run for a short time
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Each block should be processed at least twice due to retry
	mu.Lock()
	defer mu.Unlock()
	require.GreaterOrEqual(t, len(processed), 4)
	require.Equal(t, 2, errorsCount[1])
	require.Equal(t, 2, errorsCount[2])
	wg.Wait()
}

func TestBaseProducer_MonitorLoopParallel_ErrorGetLatestBlockHeight(t *testing.T) {
	mockAlloraClient := new(mocks.AlloraClientInterface)
	mockRepository := new(mocks.ProcessedBlockRepositoryInterface)
	mockService := new(mocks.ProcessorService)
	logger := zerolog.Nop()

	bm := &BaseProducer{
		alloraClient:         mockAlloraClient,
		repository:           mockRepository,
		service:              mockService,
		logger:               &logger,
		startHeight:          1,
		blockRefreshInterval: 10 * time.Millisecond,
		rateLimitInterval:    10 * time.Millisecond,
	}

	var mu sync.Mutex
	processed := []int64{}
	processBlock := func(_ context.Context, height int64) error {
		mu.Lock()
		defer mu.Unlock()
		processed = append(processed, height)
		return nil
	}

	mockAlloraClient.On("GetLatestBlockHeight", mock.Anything).Return(int64(0), errors.New("network error"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := bm.MonitorLoopParallel(ctx, processBlock, 2)
		assert.NoError(t, err)
		wg.Done()
	}()

	// Let it run for a short time
	time.Sleep(100 * time.Millisecond)
	cancel()

	require.Empty(t, processed)
	require.Equal(t, int64(1), bm.startHeight)
	wg.Wait()
}
