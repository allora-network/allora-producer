package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/allora-network/allora-producer/app/domain/mocks"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test ProcessorService.ProcessBlock
func TestProcessorService_ProcessBlock(t *testing.T) {
	mockKafka := new(mocks.StreamingClient)
	mockCodec := new(mocks.CodecInterface)
	mockEventFilter := new(mocks.FilterInterface[abcitypes.Event])
	mockTxFilter := new(mocks.FilterInterface[codectypes.Any])

	service, err := NewProcessorService(mockKafka, mockCodec, mockEventFilter, mockTxFilter)
	require.NoError(t, err)

	block := &coretypes.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height:  1,
				ChainID: "test-chain",
				Time:    time.Now(),
			},
			Data: types.Data{
				Txs: types.Txs{
					[]byte("tx1"),
					[]byte("tx2"),
				},
			},
		},
	}

	parsedTx := &txtypes.Tx{
		Body: &txtypes.TxBody{
			Messages: []*codectypes.Any{
				{
					TypeUrl: "/test.Msg",
				},
			},
		},
	}

	mockCodec.On("ParseTx", []byte("tx1")).Return(parsedTx, nil)
	mockTxFilter.On("ShouldProcess", parsedTx.Body.Messages[0]).Return(true)
	mockCodec.On("MarshalProtoJSON", parsedTx.Body.Messages[0]).Return(json.RawMessage(`{"key":"value"}`), nil)
	mockKafka.On("PublishAsync", mock.Anything, "/test.Msg", mock.Anything, mock.Anything).Return(nil)
	mockCodec.On("ParseTx", []byte("tx2")).Return(nil, errors.New("parse error"))

	err = service.ProcessBlock(context.Background(), block)
	require.NoError(t, err)

	mockKafka.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
	mockEventFilter.AssertExpectations(t)
	mockTxFilter.AssertExpectations(t)
}

// Test ProcessorService.ProcessTransaction
func TestProcessorService_ProcessTransaction(t *testing.T) {
	mockKafka := new(mocks.StreamingClient)
	mockCodec := new(mocks.CodecInterface)
	mockEventFilter := new(mocks.FilterInterface[abcitypes.Event])
	mockTxFilter := new(mocks.FilterInterface[codectypes.Any])

	service, err := NewProcessorService(mockKafka, mockCodec, mockEventFilter, mockTxFilter)
	require.NoError(t, err)

	tx := []byte("tx")
	header := &types.Header{
		Height:  1,
		ChainID: "test-chain",
		Time:    time.Now(),
	}

	parsedTx := &txtypes.Tx{
		Body: &txtypes.TxBody{
			Messages: []*codectypes.Any{
				{TypeUrl: "/test.Msg"},
			},
		},
	}

	mockCodec.On("ParseTx", tx).Return(parsedTx, nil)
	mockTxFilter.On("ShouldProcess", parsedTx.Body.Messages[0]).Return(true)
	mockCodec.On("MarshalProtoJSON", parsedTx.Body.Messages[0]).Return(json.RawMessage(`{"key":"value"}`), nil)
	mockKafka.On("PublishAsync", mock.Anything, "/test.Msg", mock.Anything, mock.Anything).Return(nil)

	err = service.ProcessTransaction(context.Background(), tx, 0, header)
	require.NoError(t, err)

	mockKafka.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
	mockTxFilter.AssertExpectations(t)
}

// Test ProcessorService.ProcessEvent
func TestProcessorService_ProcessEvent(t *testing.T) {
	mockKafka := new(mocks.StreamingClient)
	mockCodec := new(mocks.CodecInterface)
	mockEventFilter := new(mocks.FilterInterface[abcitypes.Event])
	mockTxFilter := new(mocks.FilterInterface[codectypes.Any])

	service, err := NewProcessorService(mockKafka, mockCodec, mockEventFilter, mockTxFilter)
	require.NoError(t, err)

	event := &abcitypes.Event{
		Type: "test_event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key", Value: "value", Index: false},
		},
	}
	header := &types.Header{
		Height:  1,
		ChainID: "test-chain",
		Time:    time.Now(),
	}

	parsedEvent := &codectypes.Any{TypeUrl: "/test.Event"}

	mockEventFilter.On("ShouldProcess", event).Return(true)
	mockCodec.On("IsTypedEvent", event).Return(true)
	mockCodec.On("ParseEvent", event).Return(parsedEvent, nil)
	mockCodec.On("MarshalProtoJSON", parsedEvent).Return(json.RawMessage(`{"event":"data"}`), nil)
	mockKafka.On("PublishAsync", mock.Anything, "test_event", mock.Anything, mock.Anything).Return(nil)

	err = service.ProcessEvent(context.Background(), event, header)
	require.NoError(t, err)

	mockKafka.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
	mockEventFilter.AssertExpectations(t)
}

// Test ProcessorService.ProcessBlockResults
func TestProcessorService_ProcessBlockResults(t *testing.T) {
	mockKafka := new(mocks.StreamingClient)
	mockCodec := new(mocks.CodecInterface)
	mockEventFilter := new(mocks.FilterInterface[abcitypes.Event])
	mockTxFilter := new(mocks.FilterInterface[codectypes.Any])

	service, err := NewProcessorService(mockKafka, mockCodec, mockEventFilter, mockTxFilter)
	require.NoError(t, err)

	blockResults := &coretypes.ResultBlockResults{
		Height: 1,
		TxsResults: []*abcitypes.ExecTxResult{
			{
				Events: []abcitypes.Event{
					{
						Type: "tx_event",
						Attributes: []abcitypes.EventAttribute{
							{Key: "key", Value: "value"},
						},
					},
				},
			},
		},
		FinalizeBlockEvents: []abcitypes.Event{
			{
				Type: "finalize_event",
				Attributes: []abcitypes.EventAttribute{
					{Key: "key", Value: "value"},
				},
			},
		},
	}

	header := &types.Header{
		Height:  1,
		ChainID: "test-chain",
		Time:    time.Now(),
	}

	event1 := &abcitypes.Event{
		Type: "tx_event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key", Value: "value"},
		},
	}
	event2 := &abcitypes.Event{
		Type: "finalize_event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key", Value: "value"},
		},
	}

	parsedEvent1 := &codectypes.Any{TypeUrl: "/test.TxEvent"}
	parsedEvent2 := &codectypes.Any{TypeUrl: "/test.FinalizeEvent"}

	mockEventFilter.On("ShouldProcess", event1).Return(true)
	mockCodec.On("IsTypedEvent", event1).Return(true)
	mockCodec.On("ParseEvent", event1).Return(parsedEvent1, nil)
	mockCodec.On("MarshalProtoJSON", parsedEvent1).Return(json.RawMessage(`{"event":"tx_event"}`), nil)
	mockKafka.On("PublishAsync", mock.Anything, "tx_event", mock.Anything, mock.Anything).Return(nil)

	mockEventFilter.On("ShouldProcess", event2).Return(true)
	mockCodec.On("IsTypedEvent", event2).Return(true)
	mockCodec.On("ParseEvent", event2).Return(parsedEvent2, nil)
	mockCodec.On("MarshalProtoJSON", parsedEvent2).Return(json.RawMessage(`{"event":"finalize_event"}`), nil)
	mockKafka.On("PublishAsync", mock.Anything, "finalize_event", mock.Anything, mock.Anything).Return(nil)

	err = service.ProcessBlockResults(context.Background(), blockResults, header)
	require.NoError(t, err)

	mockKafka.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
	mockEventFilter.AssertExpectations(t)
}

func TestProcessorService_NewProcessorService_Error(t *testing.T) {
	mockKafka := new(mocks.StreamingClient)
	mockCodec := new(mocks.CodecInterface)
	mockEventFilter := new(mocks.FilterInterface[abcitypes.Event])
	mockTxFilter := new(mocks.FilterInterface[codectypes.Any])

	service, err := NewProcessorService(mockKafka, mockCodec, mockEventFilter, mockTxFilter)
	require.NoError(t, err)
	require.NotNil(t, service)
}
