package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/util"
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
	defer util.LogExecutionTime(time.Now(), "GetBlockByHeight", map[string]interface{}{
		"height": height,
	})
	block, err := a.client.Block(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block %d from %s: %w", height, a.rpcURL, err)
	}

	return block, nil
}

// GetBlockResults implements domain.AlloraClientInterface.
func (a *AlloraClient) GetBlockResults(ctx context.Context, height int64) (*coretypes.ResultBlockResults, error) {
	defer util.LogExecutionTime(time.Now(), "GetBlockResults", map[string]interface{}{
		"height": height,
	})
	results, err := a.client.BlockResults(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block results for block %d from %s: %w", height, a.rpcURL, err)
	}

	return results, nil
}

// GetLatestBlockHeight implements domain.AlloraClientInterface.
func (a *AlloraClient) GetLatestBlockHeight(ctx context.Context) (int64, error) {
	defer util.LogExecutionTime(time.Now(), "GetLatestBlockHeight", map[string]interface{}{})
	status, err := a.client.Status(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve status from %s: %w", a.rpcURL, err)
	}

	if status == nil {
		return 0, fmt.Errorf("status is nil")
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

func (a *AlloraClient) GetHeader(ctx context.Context, height int64) (*coretypes.ResultHeader, error) {
	defer util.LogExecutionTime(time.Now(), "GetHeader", map[string]interface{}{
		"height": height,
	})
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
