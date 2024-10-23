package codec

import (
	"encoding/json"
	"errors"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	emissionsv4 "github.com/allora-network/allora-producer/codec/allora-chain/x/emissions/v4types"
	"github.com/allora-network/allora-producer/codec/mocks"
)

// newTestCodec creates a new Codec instance with a mock ProtoCodec for testing
func newTestCodec() (*Codec, *mocks.ProtoCodec, *mocks.SDKTypes) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	mockCodec := &mocks.ProtoCodec{}
	mockSDKTypes := &mocks.SDKTypes{}
	c := &Codec{
		interfaceRegistry: interfaceRegistry,
		codec:             mockCodec,
		sdkTypes:          mockSDKTypes,
	}
	return c, mockCodec, mockSDKTypes
}

// TestCodec_AddInterfaces tests the AddInterfaces method
func TestCodec_AddInterfaces(t *testing.T) {
	c, _, _ := newTestCodec()

	// Attempt to resolve an interface before adding it
	_, err := c.interfaceRegistry.Resolve("/emissions.v3.MsgInsertReputerPayload")
	require.Error(t, err)

	// Add a new interface
	c.AddInterfaces(emissionsv4.RegisterInterfaces)

	// Resolve the interface after adding it
	_, err = c.interfaceRegistry.Resolve("/emissions.v4.AddStakeRequest")
	require.NoError(t, err)
}

// TestCodec_ParseTx tests the ParseTx method
func TestCodec_ParseTx(t *testing.T) {
	c, mockCodec, _ := newTestCodec()

	txBytes := []byte("test tx bytes")
	expectedTx := &sdktx.Tx{}

	// Set up the mock to return expectedTx when Unmarshal is called
	mockCodec.On("Unmarshal", txBytes, mock.Anything).Run(func(args mock.Arguments) {
		txPtr, ok := args.Get(1).(*sdktx.Tx)
		require.True(t, ok)
		*txPtr = *expectedTx
	}).Return(nil)

	txMsg, err := c.ParseTx(txBytes)
	require.NoError(t, err)
	require.Equal(t, expectedTx, txMsg)

	mockCodec.AssertExpectations(t)
}

func TestCodec_ParseTx_UnmarshalError(t *testing.T) {
	c, mockCodec, _ := newTestCodec()

	txBytes := []byte("invalid tx bytes")
	unmarshalError := errors.New("unmarshal error")

	// Set up the mock to return an error when Unmarshal is called
	mockCodec.On("Unmarshal", txBytes, mock.Anything).Return(unmarshalError)

	txMsg, err := c.ParseTx(txBytes)
	require.Error(t, err)
	require.Nil(t, txMsg)
	require.Contains(t, err.Error(), "unmarshal error")

	mockCodec.AssertExpectations(t)
}

// TestCodec_ParseTxMessage tests the ParseTxMessage method
func TestCodec_ParseTxMessage(t *testing.T) {
	c, mockCodec, _ := newTestCodec()

	expectedMsg := &banktypes.MsgSend{
		FromAddress: "from",
		ToAddress:   "to",
	}

	// Create the Any message using NewAnyWithValue
	message, err := codectypes.NewAnyWithValue(expectedMsg)
	require.NoError(t, err)

	// Set up the mock to unpack expectedMsg when UnpackAny is called
	mockCodec.On("UnpackAny", message, mock.Anything).Run(func(args mock.Arguments) {
		msgPtr, ok := args.Get(1).(*sdktypes.Msg)
		require.True(t, ok)
		*msgPtr = expectedMsg
	}).Return(nil)

	msg, err := c.ParseTxMessage(message)
	require.NoError(t, err)
	require.Equal(t, expectedMsg, msg)

	mockCodec.AssertExpectations(t)
}

func TestCodec_ParseTxMessage_UnpackAnyError(t *testing.T) {
	c, mockCodec, _ := newTestCodec()

	message := &codectypes.Any{
		TypeUrl: "/invalid.type.url",
		Value:   []byte("invalid data"),
	}

	mockCodec.On("UnpackAny", message, mock.Anything).Return(errors.New("unpack error"))

	msg, err := c.ParseTxMessage(message)
	require.Error(t, err)
	require.Nil(t, msg)
	require.Contains(t, err.Error(), "unpack error")

	mockCodec.AssertExpectations(t)
}

// TestCodec_ParseTxMessages tests the ParseTxMessages method
func TestCodec_ParseTxMessages(t *testing.T) {
	c, mockCodec, _ := newTestCodec()

	// Expected messages
	expectedMsg1 := &banktypes.MsgSend{
		FromAddress: "from1",
		ToAddress:   "to1",
	}
	expectedMsg2 := &stakingtypes.MsgDelegate{
		DelegatorAddress: "delegator",
		ValidatorAddress: "validator",
	}

	// Create Any messages using NewAnyWithValue
	message1, err := codectypes.NewAnyWithValue(expectedMsg1)
	require.NoError(t, err)
	message2, err := codectypes.NewAnyWithValue(expectedMsg2)
	require.NoError(t, err)

	messages := []*codectypes.Any{message1, message2}

	// Set up the mock to unpack expected messages
	mockCodec.On("UnpackAny", message1, mock.Anything).Run(func(args mock.Arguments) {
		msgPtr, ok := args.Get(1).(*sdktypes.Msg)
		require.True(t, ok)
		*msgPtr = expectedMsg1
	}).Return(nil)

	mockCodec.On("UnpackAny", message2, mock.Anything).Run(func(args mock.Arguments) {
		msgPtr, ok := args.Get(1).(*sdktypes.Msg)
		require.True(t, ok)
		*msgPtr = expectedMsg2
	}).Return(nil)

	parsedMessages, err := c.ParseTxMessages(messages)
	require.NoError(t, err)
	require.Len(t, parsedMessages, 2)
	require.Equal(t, expectedMsg1, parsedMessages[0])
	require.Equal(t, expectedMsg2, parsedMessages[1])

	mockCodec.AssertExpectations(t)
}

// TestCodec_ParseTxMessages_WithErrors tests ParseTxMessages handling errors
func TestCodec_ParseTxMessages_WithErrors(t *testing.T) {
	c, mockCodec, _ := newTestCodec()

	// Expected messages
	expectedMsg1 := &banktypes.MsgSend{
		FromAddress: "from1",
		ToAddress:   "to1",
	}

	// Create Any messages
	message1, err := codectypes.NewAnyWithValue(expectedMsg1)
	require.NoError(t, err)

	// Second message is invalid
	message2 := &codectypes.Any{
		TypeUrl: "/invalid.type.url",
		Value:   []byte("invalid data"),
	}

	messages := []*codectypes.Any{message1, message2}

	// Set up the mock to unpack the first message and fail on the second
	mockCodec.On("UnpackAny", message1, mock.Anything).Run(func(args mock.Arguments) {
		msgPtr, ok := args.Get(1).(*sdktypes.Msg)
		require.True(t, ok)
		*msgPtr = expectedMsg1
	}).Return(nil)

	unpackError := errors.New("unpack error")
	mockCodec.On("UnpackAny", message2, mock.Anything).Return(unpackError)

	parsedMessages, err := c.ParseTxMessages(messages)
	require.Error(t, err)
	require.Len(t, parsedMessages, 1)
	require.Equal(t, expectedMsg1, parsedMessages[0])
	require.Contains(t, err.Error(), "unpack error")

	mockCodec.AssertExpectations(t)
}

// TestCodec_ParseEvent tests the ParseEvent method
func TestCodec_ParseEvent(t *testing.T) {
	c, _, mockSDKTypes := newTestCodec()

	mockEvent := &abcitypes.Event{
		Type: "test.event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}
	expectedProtoMessage := &abcitypes.Event{
		Type: "test.event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}

	mockSDKTypes.On("ParseTypedEvent", mock.Anything).Return(expectedProtoMessage, nil)

	result, err := c.ParseEvent(mockEvent)

	require.NoError(t, err)
	require.Equal(t, expectedProtoMessage, result)
	mockSDKTypes.AssertExpectations(t)
}

func TestCodec_ParseEvent_WithModeAttribute(t *testing.T) {
	c, _, mockSDKTypes := newTestCodec()

	mockEvent := &abcitypes.Event{
		Type: "test.event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "mode", Value: "value3"},
		},
	}

	mockEventCopy := *mockEvent
	mockEventCopy.Attributes = mockEventCopy.Attributes[:len(mockEventCopy.Attributes)-1]

	expectedProtoMessage := &abcitypes.Event{
		Type: "test.event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}

	mockSDKTypes.On("ParseTypedEvent", mockEventCopy).Return(expectedProtoMessage, nil)

	result, err := c.ParseEvent(mockEvent)

	require.NoError(t, err)
	require.Equal(t, expectedProtoMessage, result)
	mockSDKTypes.AssertExpectations(t)
}

// TestCodec_ParseEvent_NoAttributes tests ParseEvent with no attributes
func TestCodec_ParseEvent_NoAttributes(t *testing.T) {
	c, _, _ := newTestCodec()

	event := &abcitypes.Event{
		Type:       "test.event",
		Attributes: []abcitypes.EventAttribute{},
	}

	parsedEvent, err := c.ParseEvent(event)
	require.Error(t, err)
	require.Nil(t, parsedEvent)
	require.Equal(t, "event has no attributes", err.Error())
}

// TestCodec_MarshalProtoJSON tests the MarshalProtoJSON method
func TestCodec_MarshalProtoJSON(t *testing.T) {
	c, mockCodec, _ := newTestCodec()

	msg := &banktypes.MsgSend{
		FromAddress: "from",
		ToAddress:   "to",
	}
	expectedJSON, err := json.Marshal(msg)
	require.NoError(t, err)

	// Set up the mock to return expected JSON when MarshalJSON is called
	mockCodec.On("MarshalJSON", msg).Return(expectedJSON, nil)

	jsonMsg, err := c.MarshalProtoJSON(msg)
	require.NoError(t, err)
	require.Equal(t, json.RawMessage(expectedJSON), jsonMsg)

	mockCodec.AssertExpectations(t)
}

// TestCodec_MarshalProtoJSON_Error tests MarshalProtoJSON handling errors
func TestCodec_MarshalProtoJSON_Error(t *testing.T) {
	c, mockCodec, _ := newTestCodec()

	msg := &banktypes.MsgSend{}
	marshalError := errors.New("marshal error")

	// Set up the mock to return an error when MarshalJSON is called
	mockCodec.On("MarshalJSON", msg).Return(nil, marshalError)

	jsonMsg, err := c.MarshalProtoJSON(msg)
	require.Error(t, err)
	require.Nil(t, jsonMsg)
	require.EqualError(t, err, "failed to marshal json: marshal error")

	mockCodec.AssertExpectations(t)
}

func TestCodec_ParseEvent_WithParseTypedEventError(t *testing.T) {
	c, _, mockSDKTypes := newTestCodec()

	mockEvent := &abcitypes.Event{
		Type: "test.event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}

	parseError := errors.New("parse error")
	mockSDKTypes.On("ParseTypedEvent", mock.Anything).Return(nil, parseError)

	result, err := c.ParseEvent(mockEvent)
	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, "failed to parse typed event: parse error")
	mockSDKTypes.AssertExpectations(t)
}

func TestCodec_NewCodec(t *testing.T) {
	c := NewCodec()
	require.NotNil(t, c)
}

func TestCodec_NewWithInterfaces(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	sdkTypes := DefaultSDKTypes{}
	c := NewCodecWithInterfaces(interfaceRegistry, codec, sdkTypes)
	require.NotNil(t, c)
}

func TestCodec_ParseUntypedEvent(t *testing.T) {
	c, _, _ := newTestCodec()

	mockEvent := &abcitypes.Event{
		Type: "test.event",
		Attributes: []abcitypes.EventAttribute{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}

	expectedEvent := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	var expectedJSON json.RawMessage
	expectedJSON, err := json.Marshal(expectedEvent)
	require.NoError(t, err)

	result, err := c.ParseUntypedEvent(mockEvent)
	require.NoError(t, err)
	require.Equal(t, expectedJSON, result)
}

func TestCodec_ParseUntypedEvent_NoAttributes(t *testing.T) {
	c, _, _ := newTestCodec()

	mockEvent := &abcitypes.Event{
		Type:       "test.event",
		Attributes: []abcitypes.EventAttribute{},
	}

	result, err := c.ParseUntypedEvent(mockEvent)
	require.NoError(t, err)
	require.Equal(t, json.RawMessage("{}"), result)
}
