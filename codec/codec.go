package codec

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/allora-network/allora-producer/app/domain"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	cosmossdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"

	mint "github.com/allora-network/allora-chain/x/mint/types"

	// Make sure types are imported in this order (newest to oldest) to avoid conflicts
	emissions "github.com/allora-network/allora-chain/x/emissions/types"

	emissionsv3 "github.com/allora-network/allora-producer/codec/allora-chain/x/emissions/v3types"

	emissionsv2 "github.com/allora-network/allora-producer/codec/allora-chain/x/emissions/v2types"
)

// ProtoCodec is a wrapper around codec.ProtoCodec
//
//go:generate mockery --name=ProtoCodec
type ProtoCodec interface {
	MarshalJSON(msg proto.Message) ([]byte, error)
	Unmarshal(data []byte, msg proto.Message) error
	UnpackAny(any *codectypes.Any, iface interface{}) error
}

// Define an interface for the sdktypes package
//
//go:generate mockery --name=SDKTypes
type SDKTypes interface {
	ParseTypedEvent(event abcitypes.Event) (proto.Message, error)
}

// Default implementation of the SDKTypes interface
type DefaultSDKTypes struct{}

func (d DefaultSDKTypes) ParseTypedEvent(event abcitypes.Event) (proto.Message, error) {
	return cosmossdktypes.ParseTypedEvent(event)
}

type Codec struct {
	interfaceRegistry codectypes.InterfaceRegistry
	codec             ProtoCodec
	sdkTypes          SDKTypes
}

var _ domain.CodecInterface = &Codec{}

var defaultRegisterFuncs = []func(codectypes.InterfaceRegistry){
	std.RegisterInterfaces,
	bank.RegisterInterfaces,
	staking.RegisterInterfaces,
	slashing.RegisterInterfaces,
	distribution.RegisterInterfaces,
	// Allora types
	emissions.RegisterInterfaces,
	mint.RegisterInterfaces,
	// Allora types (local)
	emissionsv2.RegisterInterfaces,
	emissionsv3.RegisterInterfaces,
}

// NewCodec creates a new instance of Codec with registered interfaces.
func NewCodec() *Codec {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	registerInterfaces(interfaceRegistry, defaultRegisterFuncs...)
	codec := codec.NewProtoCodec(interfaceRegistry)

	return &Codec{interfaceRegistry: interfaceRegistry, codec: codec, sdkTypes: DefaultSDKTypes{}}
}

func registerInterfaces(interfaceRegistry codectypes.InterfaceRegistry, registerFuncs ...func(codectypes.InterfaceRegistry)) {
	for _, registerFunc := range registerFuncs {
		registerFunc(interfaceRegistry)
	}
}

func (c *Codec) AddInterfaces(registerFuncs ...func(codectypes.InterfaceRegistry)) {
	for _, registerFunc := range registerFuncs {
		registerFunc(c.interfaceRegistry)
	}
}

// DecodeTxMsgs decodes a transaction and returns a list of decoded messages
func (c *Codec) ParseTx(txBytes []byte) (*tx.Tx, error) {
	var txMsg tx.Tx
	err := c.codec.Unmarshal(txBytes, &txMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal txBytes: %w", err)
	}
	return &txMsg, nil
}

func (c *Codec) ParseTxMessage(message *codectypes.Any) (proto.Message, error) {
	var msg cosmossdktypes.Msg
	err := c.codec.UnpackAny(message, &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack any: %w", err)
	}
	return msg, nil
}

func (c *Codec) ParseTxMessages(messages []*codectypes.Any) ([]proto.Message, error) {
	var parsedMessages []proto.Message
	var parseErrors []error
	for _, msgAny := range messages {
		msg, err := c.ParseTxMessage(msgAny)
		if err != nil {
			parseErrors = append(parseErrors, err)
			continue
		}
		parsedMessages = append(parsedMessages, msg)
	}
	if len(parseErrors) > 0 {
		return parsedMessages, fmt.Errorf("failed to parse some messages: %v", parseErrors)
	}
	return parsedMessages, nil
}

// DecodeMsg decodes a typed event and returns a decoded event
func (c *Codec) ParseEvent(event *abcitypes.Event) (proto.Message, error) {
	// TODO: Double check this
	if len(event.Attributes) == 0 {
		return nil, errors.New("event has no attributes")
	}

	eventCopy := *event
	// Remove the last attribute if the key is "mode"
	if eventCopy.Attributes[len(eventCopy.Attributes)-1].Key == "mode" {
		eventCopy.Attributes = eventCopy.Attributes[:len(eventCopy.Attributes)-1]
	}

	protoEvent, err := c.sdkTypes.ParseTypedEvent(eventCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to parse typed event: %w", err)
	}

	return protoEvent, nil
}

// MarshalJSON marshals a proto message to a json string
func (c *Codec) MarshalProtoJSON(event proto.Message) (json.RawMessage, error) {
	jsonEvent, err := c.codec.MarshalJSON(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}
	return jsonEvent, nil
}

func (c *Codec) ParseUntypedEvent(event *abcitypes.Event) (json.RawMessage, error) {
	attrMap := make(map[string]string)
	for _, attr := range event.Attributes {
		attrMap[attr.Key] = attr.Value
	}

	attrBytes, err := json.Marshal(attrMap)
	if err != nil {
		return nil, err
	}
	return attrBytes, nil
}

func (c *Codec) IsTypedEvent(event *abcitypes.Event) bool {
	concreteGoType := proto.MessageType(event.Type)
	return concreteGoType != nil
}
