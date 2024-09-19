package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/rs/zerolog/log"
)

type App struct {
	monitorEvents       domain.MonitorEventsUseCase
	monitorTransactions domain.MonitorTransactionsUseCase
}

func NewApp(monitorEvents domain.MonitorEventsUseCase, monitorTransactions domain.MonitorTransactionsUseCase) *App {
	return &App{monitorEvents: monitorEvents, monitorTransactions: monitorTransactions}
}

func (a *App) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Start the monitorEvents use case
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("Starting monitorEvents use case")
		if err := a.monitorEvents.Execute(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start monitorEvents use case: %w", err)
		}
	}()

	// Start the monitorTransactions use case
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("Starting monitorTransactions use case")
		if err := a.monitorTransactions.Execute(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start monitorTransactions use case: %w", err)
		}
	}()

	// Wait for both use cases to complete
	wg.Wait()
	close(errChan)

	log.Info().Msg("All use cases completed")

	// Check for errors
	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errs)
	}

	return nil
}
