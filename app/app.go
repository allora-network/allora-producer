package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/rs/zerolog/log"
)

type App struct {
	eventsProducer       domain.EventsProducer
	transactionsProducer domain.TransactionsProducer
}

func NewApp(eventsProducer domain.EventsProducer, transactionsProducer domain.TransactionsProducer) *App {
	return &App{eventsProducer: eventsProducer, transactionsProducer: transactionsProducer}
}

func (a *App) Run(ctx context.Context) error {
	log.Info().Msg("Starting all producers")
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Start the eventsProducer use case
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("Starting eventsProducer use case")
		if err := a.eventsProducer.Execute(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start eventsProducer use case: %w", err)
		}
	}()

	// Start the transactionsProducer use case
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("Starting transactionsProducer use case")
		if err := a.transactionsProducer.Execute(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start transactionsProducer use case: %w", err)
		}
	}()

	// Wait for both use cases to complete
	wg.Wait()
	close(errChan)

	log.Info().Msg("All producers finished")

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
