package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/allora-network/allora-producer/infra/mocks"
)

func TestGetLastProcessedBlock_Success(t *testing.T) {
	mockPool := new(mocks.DBPool)
	repo, err := NewPgProcessedBlock(mockPool)
	require.NoError(t, err)

	mockRow := new(mocks.RowInterface)
	expectedBlock := domain.ProcessedBlock{
		ID:          1,
		Height:      123,
		ProcessedAt: time.Now(),
		Status:      "success",
	}

	mockPool.On("QueryRow", mock.Anything, "SELECT * FROM processed_blocks ORDER BY height DESC LIMIT 1").
		Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 0 {
			if blockPtr, ok := args[0].(*domain.ProcessedBlock); ok {
				*blockPtr = expectedBlock
			}
		}
	})

	block, err := repo.GetLastProcessedBlock(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedBlock, block)

	mockPool.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestGetLastProcessedBlock_NoRows(t *testing.T) {
	mockPool := new(mocks.DBPool)
	repo, err := NewPgProcessedBlock(mockPool)
	require.NoError(t, err)

	mockRow := new(mocks.RowInterface)
	mockPool.On("QueryRow", mock.Anything, "SELECT * FROM processed_blocks ORDER BY height DESC LIMIT 1").
		Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(pgx.ErrNoRows)

	block, err := repo.GetLastProcessedBlock(context.Background())
	require.NoError(t, err)
	assert.Equal(t, domain.ProcessedBlock{}, block)

	mockPool.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestGetLastProcessedBlock_Error(t *testing.T) {
	mockPool := new(mocks.DBPool)
	repo, err := NewPgProcessedBlock(mockPool)
	require.NoError(t, err)

	mockRow := new(mocks.RowInterface)
	mockPool.On("QueryRow", mock.Anything, "SELECT * FROM processed_blocks ORDER BY height DESC LIMIT 1").
		Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("scan error"))

	block, err := repo.GetLastProcessedBlock(context.Background())
	require.Error(t, err)
	assert.Equal(t, domain.ProcessedBlock{}, block)
	assert.Contains(t, err.Error(), "failed to get last processed block")

	mockPool.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestSaveProcessedBlock_Success(t *testing.T) {
	mockPool := new(mocks.DBPool)
	repo, err := NewPgProcessedBlock(mockPool)
	require.NoError(t, err)

	block := domain.ProcessedBlock{
		Height:      456,
		ProcessedAt: time.Now(),
		Status:      domain.StatusCompleted,
	}

	mockPool.On("Exec", mock.Anything, "INSERT INTO processed_blocks (height, processed_at, status) VALUES ($1, $2, $3)",
		block.Height, block.ProcessedAt, block.Status).
		Return(pgconn.CommandTag("INSERT 1"), nil)

	err = repo.SaveProcessedBlock(context.Background(), block)
	require.NoError(t, err)

	mockPool.AssertExpectations(t)
}

func TestSaveProcessedBlock_Error(t *testing.T) {
	mockPool := new(mocks.DBPool)
	repo, err := NewPgProcessedBlock(mockPool)
	require.NoError(t, err)

	block := domain.ProcessedBlock{
		Height:      789,
		ProcessedAt: time.Now(),
		Status:      "failed",
	}

	mockPool.On("Exec", mock.Anything, "INSERT INTO processed_blocks (height, processed_at, status) VALUES ($1, $2, $3)",
		block.Height, block.ProcessedAt, block.Status).
		Return(pgconn.CommandTag(""), errors.New("exec error"))

	err = repo.SaveProcessedBlock(context.Background(), block)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save processed block")

	mockPool.AssertExpectations(t)
}
