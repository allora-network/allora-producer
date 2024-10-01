package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/rs/zerolog/log"
)

type EventsProducer struct {
	BaseProducer
}

var _ domain.EventsProducer = &EventsProducer{}

func NewEventsProducer(service domain.ProcessorService, client domain.AlloraClientInterface, repository domain.ProcessedBlockRepositoryInterface,
	startHeight int64, blockRefreshInterval time.Duration, rateLimitInterval time.Duration, numWorkers int) (*EventsProducer, error) {
	if service == nil {
		return nil, fmt.Errorf("service is nil")
	}
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if repository == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	logger := log.With().Str("producer", "events").Logger()

	return &EventsProducer{
		BaseProducer: BaseProducer{
			service:              service,
			alloraClient:         client,
			repository:           repository,
			startHeight:          startHeight,
			blockRefreshInterval: blockRefreshInterval,
			rateLimitInterval:    rateLimitInterval,
			numWorkers:           numWorkers,
			logger:               &logger,
		},
	}, nil
}

func (m *EventsProducer) Execute(ctx context.Context) error {
	if err := m.InitStartHeight(ctx); err != nil {
		return err
	}

	return m.MonitorLoopParallel(ctx, m.processBlockResults, m.numWorkers)
}

func (m *EventsProducer) processBlockResults(ctx context.Context, height int64) error {
	m.logger.Info().Msgf("Processing block results for height %d", height)
	// Fetch BlockResults
	blockResults, err := m.alloraClient.GetBlockResults(ctx, height)
	if err != nil {
		return fmt.Errorf("failed to get block results for height %d: %w", height, err)
	}

	// Fetch the Header separately
	header, err := m.alloraClient.GetHeader(ctx, height)
	if err != nil {
		return fmt.Errorf("failed to get block header for height %d: %w", height, err)
	}

	// Process the block results and header
	return m.service.ProcessBlockResults(ctx, blockResults, header.Header)
}