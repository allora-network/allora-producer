package domain

import "context"

//go:generate mockery --all

type MonitorTransactionsUseCase interface {
	Execute(ctx context.Context) error
}

type MonitorEventsUseCase interface {
	Execute(ctx context.Context) error
}
