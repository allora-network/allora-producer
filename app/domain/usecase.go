package domain

import "context"

//go:generate mockery --all

type EventsProducer interface {
	Execute(ctx context.Context) error
}

type TransactionsProducer interface {
	Execute(ctx context.Context) error
}
