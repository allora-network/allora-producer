package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewTopicRouter verifies the NewTopicRouter function.
func TestNewTopicRouter(t *testing.T) {
	mapping := map[string]string{
		"TypeA": "topic_a",
		"TypeB": "topic_b",
	}

	router := NewTopicRouter(mapping)

	assert.NotNil(t, router, "Expected router to be non-nil")
	assert.Equal(t, mapping, router.mapping, "Expected mappings to be equal")
}

// TestGetTopic checks the GetTopic method of TopicRouterImpl.
func TestGetTopic(t *testing.T) {
	mapping := map[string]string{
		"TypeA": "topic_a",
		"TypeB": "topic_b",
	}
	router := NewTopicRouter(mapping)

	t.Run("ExistingMessageType", func(t *testing.T) {
		topic, err := router.GetTopic("TypeA")
		require.NoError(t, err, "Expected no error for existing message type")
		assert.Equal(t, "topic_a", topic, "Expected topic 'topic_a'")
	})

	t.Run("UnknownMessageType", func(t *testing.T) {
		expectedErr := ErrUnknownMessageType
		topic, err := router.GetTopic("TypeC")
		require.ErrorIs(t, err, expectedErr, "Expected ErrUnknownMessageType")
		assert.Empty(t, topic, "Expected topic to be empty for unknown message type")
	})
}

// TestNewTopicRouterWithNilMapping verifies the NewTopicRouter function with a nil mapping.
func TestNewTopicRouterWithNilMapping(t *testing.T) {
	router := NewTopicRouter(nil)
	assert.NotNil(t, router, "Expected router to be non-nil")
	assert.Empty(t, router.mapping, "Expected mapping to be empty")
}

// TestNewTopicRouterWithEmptyMapping verifies the NewTopicRouter function with an empty mapping.
func TestNewTopicRouterWithEmptyMapping(t *testing.T) {
	mapping := map[string]string{}
	router := NewTopicRouter(mapping)
	assert.NotNil(t, router, "Expected router to be non-nil")
	assert.Empty(t, router.mapping, "Expected router.mapping to be empty")
}

// TestGetTopicCaseSensitivity checks the GetTopic method for case sensitivity.
func TestGetTopicCaseSensitivity(t *testing.T) {
	mapping := map[string]string{
		"TypeA": "topic_a",
	}
	router := NewTopicRouter(mapping)

	t.Run("CaseSensitiveMessageType", func(t *testing.T) {
		topic, err := router.GetTopic("typea")
		require.ErrorIs(t, err, ErrUnknownMessageType, "Expected ErrUnknownMessageType")
		assert.Empty(t, topic, "Expected topic to be empty for case-sensitive message type")
	})
}
