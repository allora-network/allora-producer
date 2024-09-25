package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// DBPool defines the interface for the database pool.
// This allows us to mock the database interactions in tests.
//
//go:generate mockery --all
type DBPool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) RowInterface
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type RowInterface = pgx.Row

type pgProcessedBlock struct {
	db DBPool
}

// Ensure pgProcessedBlock implements domain.ProcessedBlockRepositoryInterface
var _ domain.ProcessedBlockRepositoryInterface = &pgProcessedBlock{}

const (
	tableProcessedBlocks      = "processed_blocks"
	tableProcessedBlockEvents = "processed_block_events"
)

// NewPgProcessedBlock creates a new instance of pgProcessedBlock
func NewPgProcessedBlock(db DBPool) (domain.ProcessedBlockRepositoryInterface, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &pgProcessedBlock{db: db}, nil
}

// GetLastProcessedBlock retrieves the last processed block from the database
func (p *pgProcessedBlock) GetLastProcessedBlock(ctx context.Context) (domain.ProcessedBlock, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY height DESC LIMIT 1", tableProcessedBlocks)
	var block domain.ProcessedBlock
	err := p.db.QueryRow(ctx, query).Scan(&block)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Handle the case where there are no rows in the table
			return domain.ProcessedBlock{}, nil
		}
		return domain.ProcessedBlock{}, fmt.Errorf("failed to get last processed block: %w", err)
	}
	return block, nil
}

// SaveProcessedBlock saves a processed block to the database
func (p *pgProcessedBlock) SaveProcessedBlock(ctx context.Context, block domain.ProcessedBlock) error {
	query := fmt.Sprintf("INSERT INTO %s (height, processed_at, status) VALUES ($1, $2, $3)", tableProcessedBlocks)
	_, err := p.db.Exec(ctx, query, block.Height, block.ProcessedAt, block.Status)
	if err != nil {
		return fmt.Errorf("failed to save processed block: %w", err)
	}
	return nil
}

// GetProcessedBlockEvent retrieves a processed block event from the database
func (p *pgProcessedBlock) GetLastProcessedBlockEvent(ctx context.Context) (domain.ProcessedBlockEvent, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY height DESC LIMIT 1", tableProcessedBlockEvents)
	var event domain.ProcessedBlockEvent
	err := p.db.QueryRow(ctx, query).Scan(&event)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ProcessedBlockEvent{}, nil
		}
		return domain.ProcessedBlockEvent{}, fmt.Errorf("failed to get processed block event: %w", err)
	}
	return event, nil
}

func (p *pgProcessedBlock) SaveProcessedBlockEvent(ctx context.Context, event domain.ProcessedBlockEvent) error {
	query := fmt.Sprintf("INSERT INTO %s (height, processed_at, status) VALUES ($1, $2, $3)", tableProcessedBlockEvents)
	_, err := p.db.Exec(ctx, query, event.Height, event.ProcessedAt, event.Status)
	if err != nil {
		return fmt.Errorf("failed to save processed block event: %w", err)
	}
	return nil
}
