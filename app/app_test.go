package app

import (
	"context"
	"errors"
	"testing"

	"github.com/allora-network/allora-producer/app/domain/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestApp_Run(t *testing.T) {
	tests := []struct {
		name                     string
		eventsProducerExecuteErr error
		transactionsExecuteErr   error
		expectedError            bool
	}{
		{
			name:                     "Both producers succeed",
			eventsProducerExecuteErr: nil,
			transactionsExecuteErr:   nil,
			expectedError:            false,
		},
		{
			name:                     "EventsProducer fails",
			eventsProducerExecuteErr: errors.New("events error"),
			transactionsExecuteErr:   nil,
			expectedError:            true,
		},
		{
			name:                     "TransactionsProducer fails",
			eventsProducerExecuteErr: nil,
			transactionsExecuteErr:   errors.New("transactions error"),
			expectedError:            true,
		},
		{
			name:                     "Both producers fail",
			eventsProducerExecuteErr: errors.New("events error"),
			transactionsExecuteErr:   errors.New("transactions error"),
			expectedError:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockEventsProducer := new(mocks.EventsProducer)
			mockTransactionsProducer := new(mocks.TransactionsProducer)

			mockEventsProducer.On("Execute", mock.Anything).Return(tt.eventsProducerExecuteErr)
			mockTransactionsProducer.On("Execute", mock.Anything).Return(tt.transactionsExecuteErr)

			app := NewApp(mockEventsProducer, mockTransactionsProducer)

			// Act
			err := app.Run(context.Background())

			// Assert
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockEventsProducer.AssertExpectations(t)
			mockTransactionsProducer.AssertExpectations(t)
		})
	}
}
