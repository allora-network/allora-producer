package infra

import (
	"context"
	"fmt"

	"github.com/allora-network/allora-producer/app/domain"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
)

//go:generate mockery --name=RPCClient
type RPCClient interface {
	Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error)
	BlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error)
	Status(ctx context.Context) (*coretypes.ResultStatus, error)
	Header(ctx context.Context, height *int64) (*coretypes.ResultHeader, error)
}

type AlloraClient struct {
	rpcURL string
	client RPCClient
}

// Ensure AlloraClient implements domain.AlloraClientInterface
var _ domain.AlloraClientInterface = &AlloraClient{}

// GetBlockByHeight implements domain.AlloraClientInterface.
func (a *AlloraClient) GetBlockByHeight(ctx context.Context, height int64) (*coretypes.ResultBlock, error) {
	block, err := a.client.Block(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block %d from %s: %w", height, a.rpcURL, err)
	}

	return block, nil
}

// GetBlockResults implements domain.AlloraClientInterface.
func (a *AlloraClient) GetBlockResults(ctx context.Context, height int64) (*coretypes.ResultBlockResults, error) {
	results, err := a.client.BlockResults(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block results for block %d from %s: %w", height, a.rpcURL, err)
	}

	return results, nil
}

// GetLatestBlockHeight implements domain.AlloraClientInterface.
func (a *AlloraClient) GetLatestBlockHeight(ctx context.Context) (int64, error) {
	status, err := a.client.Status(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve status from %s: %w", a.rpcURL, err)
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (a *AlloraClient) GetHeader(ctx context.Context, height int64) (*coretypes.ResultHeader, error) {
	header, err := a.client.Header(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve header for block %d from %s: %w", height, a.rpcURL, err)
	}

	return header, nil
}

func NewAlloraClient(rpcURL string, timeout uint) (domain.AlloraClientInterface, error) {
	rpc, err := rpchttp.NewWithTimeout(rpcURL, "/websocket", timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client from node %s: %w", rpcURL, err)
	}

	return &AlloraClient{rpcURL: rpcURL, client: rpc}, nil
}
