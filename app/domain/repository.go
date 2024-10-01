package domain

import (
	"context"
	"encoding/json"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/gogoproto/proto"
)

//go:generate mockery --all

type AlloraClientInterface interface {
	// GetLatestBlockHeight returns the latest block height
	GetLatestBlockHeight(ctx context.Context) (int64, error)
	// GetBlockByHeight returns a block by height
	GetBlockByHeight(ctx context.Context, height int64) (*ctypes.ResultBlock, error)
	// GetBlockResults returns the results of a block by height
	GetBlockResults(ctx context.Context, height int64) (*ctypes.ResultBlockResults, error)
	// GetHeader returns the header of a block by height
	GetHeader(ctx context.Context, height int64) (*ctypes.ResultHeader, error)
}

// TopicRouter defines the interface for determining Kafka topics based on message types.
type TopicRouter interface {
	// GetTopic returns the Kafka topic for a given message type
	GetTopic(msgType string) (string, error)
}

type StreamingClient interface {
	// PublishAsync publishes a message to a topic asynchronously
	PublishAsync(ctx context.Context, msgType string, message []byte, blockHeight int64) error
	// Close closes the client
	Close() error
}

// DecodeTxMsg decodes a transaction message
type CodecInterface interface {
	// ParseTx parses a transaction
	ParseTx(txBytes []byte) (*tx.Tx, error)
	// ParseTxMessage parses a transaction message
	ParseTxMessage(message *codectypes.Any) (proto.Message, error)
	// ParseTxMessages parses transaction messages
	ParseTxMessages(txMessages []*codectypes.Any) ([]proto.Message, error)
	// ParseEvent parses an event
	ParseEvent(event *abcitypes.Event) (proto.Message, error)
	// MarshalProtoJSON marshals a protobuf message to JSON
	MarshalProtoJSON(event proto.Message) (json.RawMessage, error)
}

type ProcessedBlockRepositoryInterface interface {
	GetLastProcessedBlock(ctx context.Context) (ProcessedBlock, error)
	SaveProcessedBlock(ctx context.Context, block ProcessedBlock) error

	GetLastProcessedBlockEvent(ctx context.Context) (ProcessedBlockEvent, error)
	SaveProcessedBlockEvent(ctx context.Context, blockEvent ProcessedBlockEvent) error
}
