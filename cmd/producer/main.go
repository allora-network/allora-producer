package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/allora-network/allora-producer/app"
	"github.com/allora-network/allora-producer/app/service"
	"github.com/allora-network/allora-producer/app/usecase"
	"github.com/allora-network/allora-producer/codec"
	"github.com/allora-network/allora-producer/infra"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/allora-network/allora-producer/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting Allora Chain Producers")
	// Load config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	zerolog.SetGlobalLevel(zerolog.Level(cfg.Log.Level))

	// Validate config
	if err := config.ValidateConfig(&cfg); err != nil {
		log.Fatal().Err(err).Msg("invalid config")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Info().Msgf("Received signal: %v. Shutting down...", sig)
		cancel()
	}()

	alloraClient, err := infra.NewAlloraClient(cfg.Allora.RPC, cfg.Allora.Timeout)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create allora client")
	}

	topicRouter := infra.NewTopicRouter(getTopicMapping(cfg))
	// Instantiate Franz client
	franzClient, err := infra.NewFranzClient(cfg.Kafka.Seeds, cfg.Kafka.User, cfg.Kafka.Password)
	if err != nil {
		log.Fatal().Err(err).Str("component", "KafkaClient").Msg("failed to create Kafka client")
	}
	kafkaClient, err := infra.NewKafkaClient(franzClient, topicRouter, cfg.Kafka.NumPartitions)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Kafka client")
	}
	defer kafkaClient.Close()

	codec := codec.NewCodec()

	eventFilter := service.NewFilterEvent(
		cfg.FilterEvent.Types...,
	)

	transactionMessageFilter := service.NewFilterTransactionMessage(
		cfg.FilterTransaction.Types...,
	)

	db, err := pgxpool.Connect(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create pg client")
	}
	defer db.Close()
	processedBlockRepository, err := infra.NewPgProcessedBlock(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create processed block repository")
	}

	processorService, err := service.NewProcessorService(kafkaClient, codec, eventFilter, transactionMessageFilter)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create processor service")
	}
	eventsProducer, err := usecase.NewEventsProducer(processorService, alloraClient, processedBlockRepository, 0,
		cfg.Producer.BlockRefreshInterval, cfg.Producer.RateLimitInterval, cfg.Producer.NumWorkers)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create events producer use case")
	}
	transactionsProducer, err := usecase.NewTransactionsProducer(processorService, alloraClient, processedBlockRepository, 0,
		cfg.Producer.BlockRefreshInterval, cfg.Producer.RateLimitInterval, cfg.Producer.NumWorkers)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create transactions producer use case")
	}

	app := app.NewApp(eventsProducer, transactionsProducer)

	if err := app.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to run app")
	}
}

func getTopicMapping(cfg config.Config) map[string]string {
	topicMapping := make(map[string]string)
	for _, topicRouter := range cfg.KafkaTopicRouter {
		for _, payload := range topicRouter.Types {
			topicMapping[payload] = topicRouter.Name
		}
	}
	return topicMapping
}
