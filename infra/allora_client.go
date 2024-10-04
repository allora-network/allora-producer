package infra

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/util"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"

	jsonrpc "github.com/cometbft/cometbft/rpc/jsonrpc/client"
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
	}, nil)
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
	}, nil)
	results, err := a.client.BlockResults(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block results for block %d from %s: %w", height, a.rpcURL, err)
	}

	return results, nil
}

// GetLatestBlockHeight implements domain.AlloraClientInterface.
func (a *AlloraClient) GetLatestBlockHeight(ctx context.Context) (int64, error) {
	defer util.LogExecutionTime(time.Now(), "GetLatestBlockHeight", map[string]interface{}{}, nil)
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
	}, nil)
	header, err := a.client.Header(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve header for block %d from %s: %w", height, a.rpcURL, err)
	}

	return header, nil
}

func NewAlloraClient(rpcURL string, timeout time.Duration) (domain.AlloraClientInterface, error) {
	client, err := jsonrpc.DefaultHTTPClient(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("error creating default http client")
	}

	client.Timeout = timeout
	if transport, ok := client.Transport.(*http.Transport); ok {
		transport.DisableCompression = false
		transport.ForceAttemptHTTP2 = true
		transport.MaxIdleConns = 100
		transport.IdleConnTimeout = 90 * time.Second
		transport.TLSHandshakeTimeout = 10 * time.Second
		transport.ExpectContinueTimeout = 1 * time.Second
	} else {
		return nil, fmt.Errorf("unexpected transport type: %T", client.Transport)
	}

	rpc, err := rpchttp.NewWithClient(rpcURL, "/websocket", client)

	if err != nil {
		return nil, fmt.Errorf("failed to create new client from node %s: %w", rpcURL, err)
	}

	return &AlloraClient{rpcURL: rpcURL, client: rpc}, nil
}
