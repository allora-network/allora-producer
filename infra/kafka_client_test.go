package infra

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/infra/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"
)

type MockTopicRouter struct {
	mock.Mock
}

func (m *MockTopicRouter) GetTopic(msgType string) (string, error) {
	args := m.Called(msgType)
	return args.String(0), args.Error(1)
}

func TestPublishAsync_Success(t *testing.T) {
	mockClient := new(mocks.KafkaClient)
	mockRouter := new(MockTopicRouter)

	msgType := "testType"
	message := []byte("test message")
	topic := "testTopic"

	mockRouter.On("GetTopic", msgType).Return(topic, nil)

	blockHeight := int64(100)
	record := &kgo.Record{
		Topic: topic,
		Value: message,
		Key:   []byte(fmt.Sprintf("%d", blockHeight)),
	}

	mockClient.On("Produce", mock.Anything, record, mock.Anything).Run(func(args mock.Arguments) {
		record, ok := args.Get(1).(*kgo.Record)
		require.True(t, ok)
		promise, ok := args.Get(2).(func(*kgo.Record, error))
		require.True(t, ok)
		promise(record, nil)
	}).Return()

	client, err := NewKafkaClient(mockClient, mockRouter)
	require.NoError(t, err)

	err = client.PublishAsync(context.Background(), msgType, message, blockHeight)

	require.NoError(t, err)
	mockRouter.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestPublishAsync_GetTopicError(t *testing.T) {
	mockClient := new(mocks.KafkaClient)
	mockRouter := new(MockTopicRouter)

	msgType := "unknownType"
	message := []byte("test message")
	errGetTopic := errors.New("topic not found")

	mockRouter.On("GetTopic", msgType).Return("", errGetTopic)

	client, err := NewKafkaClient(mockClient, mockRouter)
	require.NoError(t, err)

	err = client.PublishAsync(context.Background(), msgType, message, 100)

	assert.Equal(t, errGetTopic, err)
	mockRouter.AssertExpectations(t)
	mockClient.AssertNotCalled(t, "Produce", mock.Anything, mock.Anything, mock.Anything)
}

func TestPublishAsync_ProduceError(t *testing.T) {
	mockClient := new(mocks.KafkaClient)
	mockRouter := new(MockTopicRouter)

	msgType := "testType"
	message := []byte("test message")
	topic := "testTopic"
	produceErr := errors.New("produce failed")

	mockRouter.On("GetTopic", msgType).Return(topic, nil)

	mockClient.On("Produce", mock.Anything, mock.AnythingOfType("*kgo.Record"), mock.Anything).Run(func(args mock.Arguments) {
		record, ok := args.Get(1).(*kgo.Record)
		require.True(t, ok)
		promise, ok := args.Get(2).(func(*kgo.Record, error))
		require.True(t, ok)
		promise(record, produceErr)
	}).Return()

	client, err := NewKafkaClient(mockClient, mockRouter)
	require.NoError(t, err)

	err = client.PublishAsync(context.Background(), msgType, message, 100)

	require.NoError(t, err)
	mockRouter.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestClose_Success(t *testing.T) {
	mockClient := new(mocks.KafkaClient)
	mockClient.On("Flush", mock.Anything).Return(nil)
	mockClient.On("Close").Return()

	mockRouter := new(MockTopicRouter)

	client, err := NewKafkaClient(mockClient, mockRouter)
	require.NoError(t, err)

	err = client.Close()

	require.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestClose_FlushTimeout(t *testing.T) {
	mockClient := new(mocks.KafkaClient)
	flushTimeoutErr := context.DeadlineExceeded
	mockClient.On("Flush", mock.Anything).Return(flushTimeoutErr)
	mockClient.On("Close").Return()

	mockRouter := new(MockTopicRouter)

	client, err := NewKafkaClient(mockClient, mockRouter)
	require.NoError(t, err)

	err = client.Close()

	require.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestFlush(t *testing.T) {
	mockClient := new(mocks.KafkaClient)
	mockClient.On("Flush", mock.Anything).Return(nil)
	mockRouter := new(MockTopicRouter)

	client, err := NewKafkaClient(mockClient, mockRouter)
	require.NoError(t, err)

	err = client.Flush(context.Background())
	require.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestFlush_Error(t *testing.T) {
	mockClient := new(mocks.KafkaClient)
	flushErr := errors.New("flush failed")
	mockClient.On("Flush", mock.Anything).Return(flushErr)
	mockRouter := new(MockTopicRouter)

	client, err := NewKafkaClient(mockClient, mockRouter)
	require.NoError(t, err)

	err = client.Flush(context.Background())
	require.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestNewKafkaClient(t *testing.T) {
	mockClient := new(mocks.KafkaClient)
	mockRouter := new(MockTopicRouter)

	cases := []struct {
		name    string
		client  KafkaClient
		router  domain.TopicRouter
		wantErr bool
	}{
		{
			name:    "success",
			client:  mockClient,
			router:  mockRouter,
			wantErr: false,
		},
		{
			name:    "nil client",
			client:  nil,
			router:  mockRouter,
			wantErr: true,
		},
		{
			name:    "nil router",
			client:  mockClient,
			router:  nil,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewKafkaClient(tt.client, tt.router)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, client)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestNewFranzClient(t *testing.T) {
	seeds := []string{"seed1:9092", "seed2:9092"}
	user := "testUser"
	password := "testPassword"

	client, err := NewFranzClient(seeds, user, password)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.IsType(t, &kgo.Client{}, client)
}
