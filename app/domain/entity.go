package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	MessageTypeEvent       string = "event"
	MessageTypeTransaction string = "transaction"
)

// These are the messages that are sent to Kafka
type Message struct {
	ID        string    `json:"id"`        // hash of the entire payload
	Type      string    `json:"type"`      // event or transaction
	Name      string    `json:"name"`      // type of the payload (full qualified name, e.g. emisions.v3.EventNetworkLossSet)
	Timestamp time.Time `json:"timestamp"` // timestamp of this message (created at, in UTC)
	Payload   Payload   `json:"payload"`   // payload of the message
}

type Payload struct {
	Metadata Metadata        `json:"metadata"` // metadata of the block and transaction (when available)
	Data     json.RawMessage `json:"data"`     // json encoded payload (transaction or event)
}

type Metadata struct {
	BlockMetadata       BlockMetadata       `json:"block_metadata"`       // metadata of the block
	TransactionMetadata TransactionMetadata `json:"transaction_metadata"` // metadata of the transaction
}

type BlockMetadata struct {
	Height  int64     `json:"height"`   // block height
	ChainID string    `json:"chain_id"` // chain id
	Hash    string    `json:"hash"`     // block hash
	Time    time.Time `json:"time"`     // block time
}

type TransactionMetadata struct {
	TxIndex int    `json:"tx_index,omitempty"` // transaction index
	TxHash  string `json:"tx_hash,omitempty"`  // transaction hash
	Type    string `json:"type"`               // protobuf type (full qualified name, e.g. emisions.v3.EventNetworkLossSet)
}

func NewMetadata(blockHeight int64, chainID string, blockHash string, time time.Time, txIndex int, txHash string, typeURL string) Metadata {
	// Remove the leading "/" from the typeURL if it exists
	typeURL = strings.TrimPrefix(typeURL, "/")
	return Metadata{
		BlockMetadata: BlockMetadata{
			Height:  blockHeight,
			ChainID: chainID,
			Hash:    blockHash,
			Time:    time,
		},
		TransactionMetadata: TransactionMetadata{
			TxIndex: txIndex,
			TxHash:  txHash,
			Type:    typeURL,
		},
	}
}

func NewPayload(metadata Metadata, data json.RawMessage) Payload {
	return Payload{Metadata: metadata, Data: data}
}

func NewMessage(msgType, name string, payload Payload) (Message, error) {
	// Checks if the msgType is valid
	if msgType != MessageTypeEvent && msgType != MessageTypeTransaction {
		return Message{}, fmt.Errorf("invalid message type: %s", msgType)
	}

	// Checks if the name is valid
	if name == "" {
		return Message{}, fmt.Errorf("name is required")
	}

	// Calculates ID as a hash of the entire payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}
	hash := sha256.Sum256(payloadBytes)

	// Remove the leading "/" from the name if it exists
	name = strings.TrimPrefix(name, "/")

	return Message{
		ID:        "0x" + hex.EncodeToString(hash[:]),
		Type:      msgType,
		Name:      name,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}, nil
}
