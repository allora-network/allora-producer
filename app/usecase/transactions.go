package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
)

type MonitorTransactions struct {
	BaseMonitor
}

func NewMonitorTransactions(service domain.MonitorService, client domain.AlloraClientInterface, repository domain.ProcessedBlockRepositoryInterface,
	startHeight int64, blockRefreshInterval time.Duration, rateLimitInterval time.Duration) (*MonitorTransactions, error) {
	if service == nil {
		return nil, fmt.Errorf("service is nil")
	}
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if repository == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	return &MonitorTransactions{
		BaseMonitor: BaseMonitor{
			service:              service,
			client:               client,
			repository:           repository,
			startHeight:          startHeight,
			blockRefreshInterval: blockRefreshInterval,
			rateLimitInterval:    rateLimitInterval,
		},
	}, nil
}

func (m *MonitorTransactions) Execute(ctx context.Context) error {
	if err := m.InitStartHeight(ctx); err != nil {
		return err
	}

	return m.MonitorLoop(ctx, m.processBlock)
}

func (m *MonitorTransactions) processBlock(ctx context.Context, height int64) error {
	// Fetch Block
	block, err := m.client.GetBlockByHeight(ctx, height)
	if err != nil {
		return fmt.Errorf("failed to get block for height %d: %w", height, err)
	}

	// Process the block
	return m.service.ProcessBlock(ctx, block)
}
