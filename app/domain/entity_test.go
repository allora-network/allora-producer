package domain

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestNewMetadata tests the NewMetadata constructor function.
func TestNewMetadata(t *testing.T) {
	blockHeight := int64(100)
	chainID := "test-chain"
	blockHash := "abc123def456"
	currentTime := time.Now()
	txIndex := 1
	txHash := "txhash123"
	typeURL := "/emissions.v3.EventNetworkLossSet1"

	metadata := NewMetadata(blockHeight, chainID, blockHash, currentTime, txIndex, txHash, typeURL)

	if metadata.BlockMetadata.Height != blockHeight {
		t.Errorf("Expected BlockMetadata.Height to be %d, got %d", blockHeight, metadata.BlockMetadata.Height)
	}
	if metadata.BlockMetadata.ChainID != chainID {
		t.Errorf("Expected BlockMetadata.ChainID to be %s, got %s", chainID, metadata.BlockMetadata.ChainID)
	}
	if metadata.BlockMetadata.Hash != blockHash {
		t.Errorf("Expected BlockMetadata.Hash to be %s, got %s", blockHash, metadata.BlockMetadata.Hash)
	}
	if !metadata.BlockMetadata.Time.Equal(currentTime) {
		t.Errorf("Expected BlockMetadata.Time to be %v, got %v", currentTime, metadata.BlockMetadata.Time)
	}
	if metadata.TransactionMetadata.TxIndex != txIndex {
		t.Errorf("Expected TransactionMetadata.TxIndex to be %d, got %d", txIndex, metadata.TransactionMetadata.TxIndex)
	}
	if metadata.TransactionMetadata.TxHash != txHash {
		t.Errorf("Expected TransactionMetadata.TxHash to be %s, got %s", txHash, metadata.TransactionMetadata.TxHash)
	}
	expectedType := strings.TrimPrefix(typeURL, "/")
	if metadata.TransactionMetadata.Type != expectedType {
		t.Errorf("Expected TransactionMetadata.Type to be %s, got %s", expectedType, metadata.TransactionMetadata.Type)
	}
}

// TestNewPayload tests the NewPayload constructor function.
func TestNewPayload(t *testing.T) {
	metadata := Metadata{
		BlockMetadata: BlockMetadata{
			Height:  100,
			ChainID: "test-chain",
			Hash:    "abc123def456",
			Time:    time.Now(),
		},
		TransactionMetadata: TransactionMetadata{
			TxIndex: 1,
			TxHash:  "txhash123",
			Type:    "emissions.v3.EventNetworkLossSet2",
		},
	}
	data := json.RawMessage(`{"key":"value"}`)

	payload := NewPayload(metadata, data)

	if payload.Metadata != metadata {
		t.Errorf("Expected Payload.Metadata to be %v, got %v", metadata, payload.Metadata)
	}
	if string(payload.Data) != `{"key":"value"}` {
		t.Errorf("Expected Payload.Data to be %s, got %s", `{"key":"value"}`, string(payload.Data))
	}
}

// TestNewMessage tests the NewMessage constructor function.
func TestNewMessage(t *testing.T) {
	msgType := MessageTypeEvent
	name := "/emissions.v3.EventNetworkLossSet3"
	payload := Payload{
		Metadata: Metadata{
			BlockMetadata: BlockMetadata{
				Height:  100,
				ChainID: "test-chain",
				Hash:    "abc123def456",
				Time:    time.Now(),
			},
			TransactionMetadata: TransactionMetadata{
				TxIndex: 1,
				TxHash:  "txhash123",
				Type:    "emissions.v3.EventNetworkLossSet3",
			},
		},
		Data: json.RawMessage(`{"key":"value"}`),
	}

	message, err := NewMessage(msgType, name, payload)
	if err != nil {
		t.Fatalf("NewMessage returned an error: %v", err)
	}

	// Verify ID is a valid SHA256 hash with "0x" prefix
	if !strings.HasPrefix(message.ID, "0x") {
		t.Errorf("Expected ID to start with '0x', got %s", message.ID)
	}
	if len(message.ID) != 2+64 { // "0x" + 64 hex characters
		t.Errorf("Expected ID length to be 66, got %d", len(message.ID))
	}

	// Verify other fields
	if message.Type != msgType {
		t.Errorf("Expected Type to be %s, got %s", msgType, message.Type)
	}
	expectedName := "emissions.v3.EventNetworkLossSet3"
	if message.Name != expectedName {
		t.Errorf("Expected Name to be %s, got %s", expectedName, message.Name)
	}

	// Verify Timestamp is recent (within the last second)
	if time.Since(message.Timestamp) > time.Second {
		t.Errorf("Expected Timestamp to be recent, got %v", message.Timestamp)
	}

	if !reflect.DeepEqual(message.Payload, payload) {
		t.Errorf("Expected Payload to be %v, got %v", payload, message.Payload)
	}
}

func TestNewMessageInvalidType(t *testing.T) {
	msgType := "invalid"
	name := "/emissions.v3.EventNetworkLossSet4"
	payload := Payload{
		Metadata: Metadata{
			BlockMetadata: BlockMetadata{
				Height:  100,
				ChainID: "test-chain",
				Hash:    "abc123def456",
				Time:    time.Now(),
			},
			TransactionMetadata: TransactionMetadata{
				TxIndex: 1,
				TxHash:  "txhash123",
				Type:    "emissions.v3.EventNetworkLossSet4",
			},
		},
		Data: json.RawMessage(`{"key":"value"}`),
	}

	_, err := NewMessage(msgType, name, payload)
	if err == nil {
		t.Errorf("Expected NewMessage to return an error for invalid message type, got nil")
	}
}

func TestNewMessageEmptyName(t *testing.T) {
	msgType := MessageTypeEvent
	name := ""
	payload := Payload{
		Metadata: Metadata{
			BlockMetadata: BlockMetadata{
				Height:  100,
				ChainID: "test-chain",
				Hash:    "abc123def456",
				Time:    time.Now(),
			},
			TransactionMetadata: TransactionMetadata{
				TxIndex: 1,
				TxHash:  "txhash123",
				Type:    "emissions.v3.EventNetworkLossSet5",
			},
		},
		Data: json.RawMessage(`{"key":"value"}`),
	}

	_, err := NewMessage(msgType, name, payload)
	if err == nil {
		t.Errorf("Expected NewMessage to return an error for empty name, got nil")
	}
}
