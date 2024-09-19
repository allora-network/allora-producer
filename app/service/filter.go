package service

import (
	"strings"

	"github.com/allora-network/allora-producer/app/domain"
	abci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type Filter struct {
	allowedTypes map[string]bool
}

func NewFilter(allowedTypes ...string) *Filter {
	allowedTypesMap := make(map[string]bool)
	for _, t := range allowedTypes {
		allowedTypesMap[t] = true
	}
	return &Filter{allowedTypes: allowedTypesMap}
}

func (f *Filter) ShouldProcess(typeValue string) bool {
	return f.allowedTypes[typeValue]
}

// FilterEvent filters events based on the allowed types.
type FilterEvent struct {
	Filter
}

var _ domain.FilterInterface[abci.Event] = &FilterEvent{}

// NewFilterEvent creates a new FilterEvent.
func NewFilterEvent(allowedTypes ...string) *FilterEvent {
	return &FilterEvent{Filter: *NewFilter(allowedTypes...)}
}

func (f *FilterEvent) ShouldProcess(event *abci.Event) bool {
	if event == nil {
		return false
	}
	// Remove the leading "/" from the event type if it exists
	name := strings.TrimPrefix(event.Type, "/")
	return f.allowedTypes[name]
}

// FilterTransactionMessage filters transaction messages based on the allowed types.
type FilterTransactionMessage struct {
	Filter
}

var _ domain.FilterInterface[codectypes.Any] = &FilterTransactionMessage{}

// NewFilterTransactionMessage creates a new FilterTransactionMessage.
func NewFilterTransactionMessage(allowedTypes ...string) *FilterTransactionMessage {
	return &FilterTransactionMessage{Filter: *NewFilter(allowedTypes...)}
}

func (f *FilterTransactionMessage) ShouldProcess(txMsg *codectypes.Any) bool {
	if txMsg == nil {
		return false
	}
	// Remove the leading "/" from the event type if it exists
	name := strings.TrimPrefix(txMsg.TypeUrl, "/")
	return f.allowedTypes[name]
}
