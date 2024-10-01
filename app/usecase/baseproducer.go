package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/rs/zerolog"
)

type BaseProducer struct {
	service              domain.ProcessorService
	alloraClient         domain.AlloraClientInterface
	repository           domain.ProcessedBlockRepositoryInterface
	startHeight          int64
	blockRefreshInterval time.Duration
	rateLimitInterval    time.Duration
	numWorkers           int
	logger               *zerolog.Logger
}

// InitStartHeight checks if the start height is zero and fetches the latest block height.
func (bm *BaseProducer) InitStartHeight(ctx context.Context) error {
	if bm.startHeight == 0 {
		latestHeight, err := bm.alloraClient.GetLatestBlockHeight(ctx)
		if err != nil {
			return fmt.Errorf("failed to get latest block height: %w", err)
		}
		bm.startHeight = latestHeight
	}
	return nil
}

// MonitorLoop allows each specific producer to define how it processes blocks or block results.
func (bm *BaseProducer) MonitorLoop(ctx context.Context, processBlock func(ctx context.Context, height int64) error) error {
	for {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Process the block or block results at the current height
		if err := processBlock(ctx, bm.startHeight); err != nil {
			bm.logger.Warn().Err(err).Msgf("failed to process block at height %d", bm.startHeight)
			time.Sleep(bm.rateLimitInterval)
			continue
		}

		// Increment the block height
		bm.startHeight++

		// Sleep between iterations to avoid spamming the node
		time.Sleep(bm.blockRefreshInterval)
	}
}

func (bm *BaseProducer) MonitorLoopParallel(ctx context.Context, processBlock func(ctx context.Context, height int64) error, numWorkers int) error {
	blockQueue := make(chan int64, 100) // Buffered channel for block heights

	// Producer Goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(blockQueue)
				return
			default:
			}

			latestHeight, err := bm.alloraClient.GetLatestBlockHeight(ctx)
			if err != nil {
				bm.logger.Error().Err(err).Msg("failed to get latest block height")
				time.Sleep(bm.rateLimitInterval)
				continue
			}

			for bm.startHeight <= latestHeight {
				blockQueue <- bm.startHeight
				bm.startHeight++
			}

			time.Sleep(bm.blockRefreshInterval)
		}
	}()

	// Consumer Goroutines
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for height := range blockQueue {
				if err := processBlock(ctx, height); err != nil {
					bm.logger.Warn().Err(err).Msgf("failed to process block at height %d", height)
					// Re-enqueue for immediate retry
					blockQueue <- height
				}
				time.Sleep(bm.rateLimitInterval)
			}
		}()
	}

	wg.Wait()
	return nil
}