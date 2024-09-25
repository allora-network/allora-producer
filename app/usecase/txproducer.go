package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
)

type TransactionsProducer struct {
	BaseProducer
}

var _ domain.TransactionsProducer = &TransactionsProducer{}

func NewTransactionsProducer(service domain.ProcessorService, client domain.AlloraClientInterface, repository domain.ProcessedBlockRepositoryInterface,
	startHeight int64, blockRefreshInterval time.Duration, rateLimitInterval time.Duration) (*TransactionsProducer, error) {
	if service == nil {
		return nil, fmt.Errorf("service is nil")
	}
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if repository == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	return &TransactionsProducer{
		BaseProducer: BaseProducer{
			service:              service,
			client:               client,
			repository:           repository,
			startHeight:          startHeight,
			blockRefreshInterval: blockRefreshInterval,
			rateLimitInterval:    rateLimitInterval,
		},
	}, nil
}

func (m *TransactionsProducer) Execute(ctx context.Context) error {
	if err := m.InitStartHeight(ctx); err != nil {
		return err
	}

	return m.MonitorLoop(ctx, m.processBlock)
}

func (m *TransactionsProducer) processBlock(ctx context.Context, height int64) error {
	// Fetch Block
	block, err := m.client.GetBlockByHeight(ctx, height)
	if err != nil {
		return fmt.Errorf("failed to get block for height %d: %w", height, err)
	}

	// Process the block
	return m.service.ProcessBlock(ctx, block)
}
