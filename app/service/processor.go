package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/rs/zerolog/log"

	"github.com/allora-network/allora-producer/app/domain"
	abci "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
)

var (
	ErrParseTransaction = errors.New("failed to parse transaction")
)

type ProcessorService struct {
	streamingClient          domain.StreamingClient
	codec                    domain.CodecInterface
	filterEvent              domain.FilterInterface[abci.Event]
	filterTransactionMessage domain.FilterInterface[codectypes.Any]
}

var _ domain.ProcessorService = &ProcessorService{}

func NewProcessorService(streamingClient domain.StreamingClient, codec domain.CodecInterface, filterEvent domain.FilterInterface[abci.Event],
	filterTransactionMessage domain.FilterInterface[codectypes.Any]) (*ProcessorService, error) {
	if streamingClient == nil {
		return nil, fmt.Errorf("kafkaClient is nil")
	}
	if codec == nil {
		return nil, fmt.Errorf("codec is nil")
	}
	if filterEvent == nil {
		return nil, fmt.Errorf("filterEvent is nil")
	}
	if filterTransactionMessage == nil {
		return nil, fmt.Errorf("filterTransactionMessage is nil")
	}

	return &ProcessorService{streamingClient: streamingClient, codec: codec, filterEvent: filterEvent, filterTransactionMessage: filterTransactionMessage}, nil
}

// ProcessBlock implements domain.ProcessorService.
func (m *ProcessorService) ProcessBlock(ctx context.Context, block *coretypes.ResultBlock) error {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	if block.Block == nil {
		return fmt.Errorf("block.Block is nil")
	}

	// defer util.LogExecutionTime(time.Now(), "ProcessBlock", map[string]interface{}{
	// 	"height": block.Block.Header.Height,
	// 	"txs":    len(block.Block.Data.Txs),
	// })

	for i, tx := range block.Block.Data.Txs {
		err := m.ProcessTransaction(ctx, tx, i, &block.Block.Header)
		if err != nil {
			if errors.Is(err, ErrParseTransaction) {
				// Probably a missing codec
				// Log and continue
				log.Warn().Err(err).Msgf("failed to parse transaction %d, height %d", i, block.Block.Header.Height)
			} else {
				return fmt.Errorf("failed to process transaction %d: %w", i, err)
			}
		}
	}

	log.Info().Int64("height", block.Block.Header.Height).Msg("block processed successfully")
	return nil
}

func (m *ProcessorService) ProcessTransaction(ctx context.Context, tx []byte, txIndex int, header *types.Header) error {
	if header == nil {
		return fmt.Errorf("header is nil")
	}

	parsedTx, err := m.codec.ParseTx(tx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse transaction")
		return ErrParseTransaction
	}

	for _, txMsg := range parsedTx.Body.Messages {
		// Ensure that only relevant messages are processed
		if !m.filterTransactionMessage.ShouldProcess(txMsg) {
			log.Debug().Str("typeUrl", txMsg.TypeUrl).Msg("transaction message filtered out")
			continue
		}

		jsonMsg, err := m.codec.MarshalProtoJSON(txMsg)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction message: %w", err)
		}

		txHash := strings.ToUpper(hex.EncodeToString(types.Tx(tx).Hash()))
		metadata := domain.NewMetadata(header.Height, header.ChainID, header.Hash().String(), header.Time, txIndex, txHash, txMsg.TypeUrl)
		payload := domain.NewPayload(metadata, jsonMsg)
		message, err := domain.NewMessage(domain.MessageTypeTransaction, txMsg.TypeUrl, payload)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		err = m.streamingClient.PublishAsync(ctx, txMsg.TypeUrl, jsonMessage, header.Height)
		if err != nil {
			return fmt.Errorf("failed to publish transaction message: %w", err)
		}

		log.Debug().Str("typeUrl", txMsg.TypeUrl).Str("message", string(jsonMessage)).Msg("transaction message processed successfully")
	}
	return nil
}

// ProcessBlockResults implements domain.ProcessorService.
func (m *ProcessorService) ProcessBlockResults(ctx context.Context, blockResults *coretypes.ResultBlockResults, header *types.Header) error {
	if blockResults == nil {
		return fmt.Errorf("blockResults is nil")
	}

	if header == nil {
		return fmt.Errorf("header is nil")
	}

	// defer util.LogExecutionTime(time.Now(), "ProcessBlockResults", map[string]interface{}{
	// 	"height":                header.Height,
	// 	"tx results":            len(blockResults.TxsResults),
	// 	"finalize block events": len(blockResults.FinalizeBlockEvents),
	// })

	for i, txResult := range blockResults.TxsResults {
		if txResult == nil {
			log.Warn().Int64("height", blockResults.Height).Int("txIndex", i).Msg("transaction result is nil")
			continue
		}

		for _, event := range txResult.Events {
			err := m.ProcessEvent(ctx, &event, header)
			if err != nil {
				log.Warn().Err(err).Int64("height", blockResults.Height).Int("txIndex", i).Msg("failed to process event")
				continue
			}
		}
	}

	for _, event := range blockResults.FinalizeBlockEvents {
		err := m.ProcessEvent(ctx, &event, header)
		if err != nil {
			log.Warn().Err(err).Int64("height", blockResults.Height).Msg("failed to process event")
			continue
		}
	}

	log.Info().Int64("height", blockResults.Height).Msg("block results processed successfully")
	return nil
}

func (m *ProcessorService) ProcessEvent(ctx context.Context, event *abci.Event, header *types.Header) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	if header == nil {
		return fmt.Errorf("header is nil")
	}

	// Ensure that only relevant events are processed
	if !m.filterEvent.ShouldProcess(event) {
		return nil
	}

	msg, err := m.codec.ParseEvent(event)
	if err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}
	jsonMsg, err := m.codec.MarshalProtoJSON(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	metadata := domain.NewMetadata(header.Height, header.ChainID, header.Hash().String(), header.Time, 0, "", event.Type)
	payload := domain.NewPayload(metadata, jsonMsg)
	message, err := domain.NewMessage(domain.MessageTypeEvent, event.Type, payload)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = m.streamingClient.PublishAsync(ctx, event.Type, jsonMessage, header.Height)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Debug().Str("eventType", event.Type).Str("event", string(jsonMessage)).Msg("event processed successfully")
	return nil
}
