package service

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/assert"
)

// TestNewFilter tests the NewFilter constructor and ShouldProcess method.
func TestNewFilter(t *testing.T) {
	allowedTypes := []string{"type1", "type2"}
	filter := NewFilter(allowedTypes...)

	for _, tpe := range allowedTypes {
		assert.True(t, filter.ShouldProcess(tpe), "Expected type %s to be allowed", tpe)
	}

	assert.False(t, filter.ShouldProcess("type3"), "Expected type 'type3' to be not allowed")
}

// TestNewFilterEvent tests the NewFilterEvent constructor and ShouldProcess method for FilterEvent.
func TestNewFilterEvent(t *testing.T) {
	allowedTypes := []string{"event1", "event2"}
	filterEvent := NewFilterEvent(allowedTypes...)

	for _, tpe := range allowedTypes {
		event := &abci.Event{Type: tpe}
		assert.True(t, filterEvent.ShouldProcess(event), "Expected event type %s to be allowed", tpe)
	}

	event := &abci.Event{Type: "event3"}
	assert.False(t, filterEvent.ShouldProcess(event), "Expected event type 'event3' to be not allowed")
}

// TestNewFilterTransactionMessage tests the NewFilterTransactionMessage constructor and ShouldProcess method for FilterTransactionMessage.
func TestNewFilterTransactionMessage(t *testing.T) {
	allowedTypes := []string{"msg1", "msg2"}
	filterTxMsg := NewFilterTransactionMessage(allowedTypes...)

	for _, tpe := range allowedTypes {
		txMsg := &codectypes.Any{TypeUrl: tpe}
		assert.True(t, filterTxMsg.ShouldProcess(txMsg), "Expected message type %s to be allowed", tpe)
	}

	txMsg := &codectypes.Any{TypeUrl: "msg3"}
	assert.False(t, filterTxMsg.ShouldProcess(txMsg), "Expected message type 'msg3' to be not allowed")
}

// TestEmptyAllowedTypes tests the behavior when the allowed types list is empty.
func TestEmptyAllowedTypes(t *testing.T) {
	filter := NewFilter()
	assert.False(t, filter.ShouldProcess("type1"), "Expected no types to be allowed")
}

// TestNilEvent tests the behavior when a nil event is passed to ShouldProcess.
func TestNilEvent(t *testing.T) {
	filterEvent := NewFilterEvent("event1")
	assert.False(t, filterEvent.ShouldProcess(nil), "Expected nil event to be not allowed")
}

// TestNilTransactionMessage tests the behavior when a nil transaction message is passed to ShouldProcess.
func TestNilTransactionMessage(t *testing.T) {
	filterTxMsg := NewFilterTransactionMessage("msg1")
	assert.False(t, filterTxMsg.ShouldProcess(nil), "Expected nil message to be not allowed")
}

// TestCaseSensitivity tests if the filter is case-sensitive.
func TestCaseSensitivity(t *testing.T) {
	filter := NewFilter("Type1")
	assert.False(t, filter.ShouldProcess("type1"), "Expected type 'type1' to be not allowed due to case sensitivity")
}

// TestDuplicateTypes tests the behavior when duplicate types are included in the allowed types list.
func TestDuplicateTypes(t *testing.T) {
	filter := NewFilter("type1", "type1")
	assert.True(t, filter.ShouldProcess("type1"), "Expected type 'type1' to be allowed even if duplicated")
}
