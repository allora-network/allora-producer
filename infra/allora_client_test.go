package infra

import (
	"context"
	"errors"
	"testing"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/allora-network/allora-producer/infra/mocks"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
)

func TestAlloraClient_GetBlockByHeight(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()
	height := int64(100)
	expectedBlock := &coretypes.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height: height,
			},
			Data: types.Data{
				Txs: types.Txs{
					[]byte("tx1"),
					[]byte("tx2"),
				},
			},
		},
	}

	mockRPC.On("Block", ctx, &height).Return(expectedBlock, nil)

	block, err := client.GetBlockByHeight(ctx, height)
	require.NoError(t, err)
	assert.Equal(t, expectedBlock, block)
	mockRPC.AssertExpectations(t)
}

func TestAlloraClient_GetBlockByHeight_Error(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()
	height := int64(100)
	expectedError := errors.New("rpc error")

	mockRPC.On("Block", ctx, &height).Return(nil, expectedError)

	block, err := client.GetBlockByHeight(ctx, height)
	require.Error(t, err)
	assert.Nil(t, block)
	assert.Contains(t, err.Error(), "failed to retrieve block")
	mockRPC.AssertExpectations(t)
}

func TestAlloraClient_GetBlockResults(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()
	height := int64(100)
	expectedResults := &coretypes.ResultBlockResults{
		Height: height,
		TxsResults: []*abci.ExecTxResult{
			{
				Code: 0,
				Data: []byte("result data"),
				Events: []abci.Event{
					{
						Type: "event_type",
						Attributes: []abci.EventAttribute{
							{
								Key:   "key1",
								Value: "value1",
								Index: true,
							},
						},
					},
				},
			},
		},
		FinalizeBlockEvents: []abci.Event{
			{
				Type: "finalize_event_type",
				Attributes: []abci.EventAttribute{
					{
						Key:   "key2",
						Value: "value2",
						Index: true,
					},
				},
			},
		},
	}

	mockRPC.On("BlockResults", ctx, &height).Return(expectedResults, nil)

	results, err := client.GetBlockResults(ctx, height)
	require.NoError(t, err)
	assert.Equal(t, expectedResults, results)
	mockRPC.AssertExpectations(t)
}

func TestAlloraClient_GetBlockResults_Error(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()
	height := int64(100)
	expectedError := errors.New("rpc error")

	mockRPC.On("BlockResults", ctx, &height).Return(nil, expectedError)

	results, err := client.GetBlockResults(ctx, height)
	require.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "failed to retrieve block results")
	mockRPC.AssertExpectations(t)
}

func TestAlloraClient_GetLatestBlockHeight(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()
	expectedHeight := int64(150)
	expectedStatus := &coretypes.ResultStatus{
		SyncInfo: coretypes.SyncInfo{
			LatestBlockHeight: expectedHeight,
		},
	}

	mockRPC.On("Status", ctx).Return(expectedStatus, nil)

	height, err := client.GetLatestBlockHeight(ctx)
	require.NoError(t, err)
	assert.Equal(t, expectedHeight, height)
	mockRPC.AssertExpectations(t)
}

func TestAlloraClient_GetLatestBlockHeight_Error(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()
	expectedError := errors.New("rpc error")

	mockRPC.On("Status", ctx).Return(nil, expectedError)

	height, err := client.GetLatestBlockHeight(ctx)
	require.Error(t, err)
	assert.Equal(t, int64(0), height)
	assert.Contains(t, err.Error(), "failed to retrieve status")
	mockRPC.AssertExpectations(t)
}

func TestAlloraClient_GetLatestBlockHeight_NilStatus(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()

	mockRPC.On("Status", ctx).Return((*coretypes.ResultStatus)(nil), nil)

	height, err := client.GetLatestBlockHeight(ctx)
	require.Error(t, err)
	assert.Equal(t, int64(0), height)
	assert.Contains(t, err.Error(), "status is nil")
	mockRPC.AssertExpectations(t)
}

func TestAlloraClient_GetHeader(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()
	height := int64(200)
	expectedHeader := &coretypes.ResultHeader{
		Header: &types.Header{
			Height: height,
		},
	}

	mockRPC.On("Header", ctx, &height).Return(expectedHeader, nil)

	header, err := client.GetHeader(ctx, height)
	require.NoError(t, err)
	assert.Equal(t, expectedHeader, header)
	mockRPC.AssertExpectations(t)
}

func TestAlloraClient_GetHeader_Error(t *testing.T) {
	mockRPC := new(mocks.RPCClient)
	client := &AlloraClient{
		rpcURL: "http://localhost:26657",
		client: mockRPC,
	}

	ctx := context.Background()
	height := int64(200)
	expectedError := errors.New("rpc error")

	mockRPC.On("Header", ctx, &height).Return(nil, expectedError)

	header, err := client.GetHeader(ctx, height)
	require.Error(t, err)
	assert.Nil(t, header)
	assert.Contains(t, err.Error(), "failed to retrieve header")
	mockRPC.AssertExpectations(t)
}

func TestNewAlloraClient(t *testing.T) {
	rpcURL := "http://localhost:26657"
	timeout := uint(10)

	clientInterface, err := NewAlloraClient(rpcURL, timeout)
	require.NoError(t, err)
	assert.NotNil(t, clientInterface)
}

func TestNewAlloraClient_Error(t *testing.T) {
	// Assuming rpchttp.NewWithTimeout returns an error for invalid URL
	invalidRPCURL := "http://invalid:timeout"
	timeout := uint(10)

	clientInterface, err := NewAlloraClient(invalidRPCURL, timeout)
	require.Error(t, err)
	assert.Nil(t, clientInterface)
	assert.Contains(t, err.Error(), "failed to create new client from node")
}
