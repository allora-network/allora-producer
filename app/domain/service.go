package domain

import (
	"context"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
)

//go:generate mockery --all

type MonitorService interface {
	ProcessBlock(ctx context.Context, block *ctypes.ResultBlock) error
	ProcessBlockResults(ctx context.Context, blockResults *ctypes.ResultBlockResults, header *types.Header) error
}

// FilterInterface defines the interface for filters.
type FilterInterface[T any] interface {
	ShouldProcess(typeValue *T) bool
}
