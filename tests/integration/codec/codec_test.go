package codec

import (
	"encoding/json"
	"testing"

	"github.com/allora-network/allora-producer/codec"
	proto "github.com/cosmos/gogoproto/proto"

	emissions "github.com/allora-network/allora-chain/x/emissions/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cosmoscodec "github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmossdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/stretchr/testify/suite"

	allora_chain_math "github.com/allora-network/allora-chain/math"

	cosmossdkmath "cosmossdk.io/math"
)

type CodecTestSuite struct {
	suite.Suite
	protoCodec *cosmoscodec.ProtoCodec
	codec      *codec.Codec
	sdkTypes   codec.SDKTypes
}

func (s *CodecTestSuite) SetupSuite() {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	// Register Cosmos SDK interfaces
	banktypes.RegisterInterfaces(interfaceRegistry)
	emissions.RegisterInterfaces(interfaceRegistry)
	s.protoCodec = cosmoscodec.NewProtoCodec(interfaceRegistry)
	s.sdkTypes = codec.DefaultSDKTypes{}
	s.codec = codec.NewCodecWithInterfaces(interfaceRegistry, s.protoCodec, s.sdkTypes)
}

func getAbciEvent(msg proto.Message) (*abci.Event, error) {
	event, err := cosmossdktypes.TypedEventToEvent(msg)
	if err != nil {
		return nil, err
	}
	abciEvent := abci.Event(event)
	return &abciEvent, nil
}

func getTx(msgs []proto.Message) (*tx.Tx, error) {
	var anyMsgs []*codectypes.Any
	for _, msg := range msgs {
		msgAny, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
		anyMsgs = append(anyMsgs, msgAny)
	}
	return &tx.Tx{
		Body: &tx.TxBody{
			Messages: anyMsgs,
		},
	}, nil
}

// TestParseAndMarshalTx tests parsing and marshaling of a transaction.
func (s *CodecTestSuite) TestParseAndMarshalTx() {
	msgSend := &banktypes.MsgSend{
		FromAddress: "cosmos1senderaddress",
		ToAddress:   "cosmos1recipientaddress",
		Amount:      cosmossdktypes.Coins{{Denom: "allo", Amount: cosmossdkmath.NewInt(100)}},
	}
	expectedJSON := "{\"@type\":\"/cosmos.bank.v1beta1.MsgSend\",\"from_address\":\"cosmos1senderaddress\",\"to_address\":\"cosmos1recipientaddress\",\"amount\":[{\"denom\":\"allo\",\"amount\":\"100\"}]}"

	sampleTx, err := getTx([]proto.Message{msgSend})
	s.Require().NoError(err, "Failed to create transaction")

	// Marshal the transaction to bytes
	txBytes, err := s.protoCodec.Marshal(sampleTx)
	s.Require().NoError(err, "Failed to marshal transaction")

	// Parse the transaction bytes
	parsedTx, err := s.codec.ParseTx(txBytes)
	s.Require().NoError(err, "Failed to parse transaction")
	s.Require().NotNil(parsedTx, "Parsed transaction is nil")
	s.Require().Len(parsedTx.Body.Messages, 1, "Expected 1 message in the transaction")

	msg := parsedTx.Body.Messages[0]
	s.Require().Equal(sampleTx.Body.Messages[0].TypeUrl, msg.TypeUrl, "Message TypeUrl does not match")
	jsonMsg, err := s.codec.MarshalProtoJSON(msg)
	s.Require().NoError(err, "Failed to marshal message to JSON")
	s.Require().Equal(expectedJSON, string(jsonMsg), "JSON does not match")
}

// TestParseAndMarshalEvent tests parsing and marshaling of an event.
func (s *CodecTestSuite) TestParseAndMarshalEvent() {
	sampleEvent := &emissions.EventRewardsSettled{
		ActorType:   1, // ACTOR_TYPE_FORECASTER
		TopicId:     2,
		BlockHeight: 3,
		Addresses:   []string{"allora1recipientaddress"},
		Rewards:     []allora_chain_math.Dec{allora_chain_math.OneDec()},
	}
	expectedJSON := "{\"actor_type\":\"ACTOR_TYPE_FORECASTER\",\"topic_id\":\"2\",\"block_height\":\"3\",\"addresses\":[\"allora1recipientaddress\"],\"rewards\":[\"1\"]}"

	abciEvent, err := getAbciEvent(sampleEvent)
	s.Require().NoError(err, "Failed to get abci event")
	// Parse the event
	parsedEvent, err := s.codec.ParseEvent(abciEvent)
	s.Require().NoError(err, "Failed to parse event")
	s.Require().NotNil(parsedEvent, "Parsed event is nil")

	// Marshal the parsed event to JSON
	jsonBytes, err := s.codec.MarshalProtoJSON(parsedEvent)
	s.Require().NoError(err, "Failed to marshal event to JSON")

	s.Require().Equal(expectedJSON, string(jsonBytes), "JSON does not match")
}

// TestIsTypedEvent tests the IsTypedEvent method.
func (s *CodecTestSuite) TestIsTypedEvent() {
	// Create a typed event
	sample := &emissions.EventRewardsSettled{}
	typedEvent, err := getAbciEvent(sample)
	s.Require().NoError(err, "Failed to get abci event")

	// Create an untyped event
	untypedEvent := &abci.Event{
		Type: "untyped_event",
	}

	// Check if the events are typed
	isTyped := s.codec.IsTypedEvent(typedEvent)
	s.Require().True(isTyped, "Expected typedEvent to be recognized as a typed event")

	isUntyped := s.codec.IsTypedEvent(untypedEvent)
	s.Require().False(isUntyped, "Expected untypedEvent to be recognized as untyped")
}

// TestParseTxMessages tests parsing multiple transaction messages.
func (s *CodecTestSuite) TestParseTxMessages() {
	// Create sample messages
	msg1 := &banktypes.MsgSend{
		FromAddress: "cosmos1senderaddress1",
		ToAddress:   "cosmos1recipientaddress1",
		Amount:      cosmossdktypes.Coins{{Denom: "allo", Amount: cosmossdkmath.NewInt(100)}},
	}

	msg2 := &banktypes.MsgSend{
		FromAddress: "cosmos1senderaddress2",
		ToAddress:   "cosmos1recipientaddress2",
		Amount:      cosmossdktypes.Coins{{Denom: "allo", Amount: cosmossdkmath.NewInt(200)}},
	}

	sampleTx, err := getTx([]proto.Message{msg1, msg2})
	s.Require().NoError(err, "Failed to create transaction")

	// Marshal the transaction to bytes
	txBytes, err := s.protoCodec.Marshal(sampleTx)
	s.Require().NoError(err, "Failed to marshal transaction")

	// Parse the transaction bytes
	parsedTx, err := s.codec.ParseTx(txBytes)
	s.Require().NoError(err, "Failed to parse transaction")
	s.Require().Len(parsedTx.Body.Messages, 2, "Expected two messages in the transaction")

	// Parse individual messages
	parsedMessages, err := s.codec.ParseTxMessages(parsedTx.Body.Messages)
	s.Require().NoError(err, "Failed to parse messages")
	s.Require().Len(parsedMessages, 2, "Expected two messages in the transaction")
}

// TestParseUntypedEvent tests parsing an untyped event.
func (s *CodecTestSuite) TestParseUntypedEvent() {
	// Create an untyped event
	untypedEvent := &abci.Event{
		Type: "custom_event",
		Attributes: []abci.EventAttribute{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}

	// Parse the untyped event
	jsonBytes, err := s.codec.ParseUntypedEvent(untypedEvent)
	s.Require().NoError(err, "Failed to parse untyped event")

	// Unmarshal JSON to map
	var eventMap map[string]string
	err = json.Unmarshal(jsonBytes, &eventMap)
	s.Require().NoError(err, "Failed to unmarshal JSON untyped event")

	// Verify the contents
	s.Require().Equal("value1", eventMap["key1"], "key1 value mismatch")
	s.Require().Equal("value2", eventMap["key2"], "key2 value mismatch")
}

func TestCodec(t *testing.T) {
	suite.Run(t, new(CodecTestSuite))
}
