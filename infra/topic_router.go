package infra

import (
	"errors"
	"strings"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/rs/zerolog/log"
)

// ErrUnknownMessageType is returned when the message type is not recognized.
var ErrUnknownMessageType = errors.New("unknown message type")

// TopicRouterImpl determines the Kafka topic for a given message type.
type TopicRouterImpl struct {
	mapping map[string]string
}

// Ensure TopicRouter implements domain.TopicRouter
var _ domain.TopicRouter = &TopicRouterImpl{}

// NewTopicRouter initializes a new TopicRouter with the provided mappings.
func NewTopicRouter(mapping map[string]string) *TopicRouterImpl {
	if mapping == nil {
		return &TopicRouterImpl{
			mapping: make(map[string]string),
		}
	}
	return &TopicRouterImpl{
		mapping: mapping,
	}
}

// GetTopic returns the Kafka topic for the given message type.
func (tr *TopicRouterImpl) GetTopic(msgType string) (string, error) {
	// Remove the leading "/" from the msgType if it exists
	name := strings.TrimPrefix(msgType, "/")
	topic, exists := tr.mapping[name]
	if !exists {
		log.Warn().Str("msgType", msgType).Msg("unknown message type")
		return "", ErrUnknownMessageType
	}
	return topic, nil
}
