package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/rs/zerolog/log"
)

type BaseProducer struct {
	service              domain.ProcessorService
	client               domain.AlloraClientInterface
	repository           domain.ProcessedBlockRepositoryInterface
	startHeight          int64
	blockRefreshInterval time.Duration
	rateLimitInterval    time.Duration
}

// InitStartHeight checks if the start height is zero and fetches the latest block height.
func (bm *BaseProducer) InitStartHeight(ctx context.Context) error {
	if bm.startHeight == 0 {
		latestHeight, err := bm.client.GetLatestBlockHeight(ctx)
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
			log.Warn().Err(err).Msgf("failed to process block at height %d", bm.startHeight)
			time.Sleep(bm.rateLimitInterval)
			continue
		}

		// Increment the block height
		bm.startHeight++

		// Sleep between iterations to avoid spamming the node
		time.Sleep(bm.blockRefreshInterval)
	}
}
