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
	"github.com/allora-network/allora-producer/util"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/allora-network/allora-producer/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Str("revision", util.Revision()).Msg("Allora Producers -- Starting")
	// Load config
	cfg, err := initializeConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Allora Producers -- Error loading config")
	}

	initializeLogging(cfg.Log.Level)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	setupSignalHandler(cancel)

	appInstance, cleanup, err := initializeApp(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Allora Producers -- Error initializing app")
	}
	defer cleanup()

	if err := appInstance.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Allora Producers -- Error running app")
	}
	log.Info().Msg("Allora Producers -- Finished")
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

func initializeConfig() (config.Config, error) {
	cfgProvider := config.NewProvider(config.NewViperAdapter())
	cfg, err := cfgProvider.InitConfig()
	if err != nil {
		return cfg, err
	}

	if err := config.ValidateConfig(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func initializeLogging(level int8) {
	zerolog.SetGlobalLevel(zerolog.Level(level))
}

func setupSignalHandler(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Info().Msgf("Received signal: %v. Shutting down...", sig)
		cancel()
	}()
}

func initializeApp(ctx context.Context, cfg config.Config) (*app.App, func(), error) {
	alloraClient, err := infra.NewAlloraClient(cfg.Allora.RPC, cfg.Allora.Timeout)
	if err != nil {
		return nil, nil, err
	}

	topicRouter := infra.NewTopicRouter(getTopicMapping(cfg))

	franzClient, err := infra.NewFranzClient(cfg.Kafka.Seeds, cfg.Kafka.User, cfg.Kafka.Password)
	if err != nil {
		return nil, nil, err
	}

	kafkaClient, err := infra.NewKafkaClient(franzClient, topicRouter)
	if err != nil {
		return nil, nil, err
	}

	codec := codec.NewCodec()

	eventFilter := service.NewFilterEvent(cfg.FilterEvent.Types...)
	transactionMessageFilter := service.NewFilterTransactionMessage(cfg.FilterTransaction.Types...)

	db, err := pgxpool.Connect(ctx, cfg.Database.URL)
	if err != nil {
		return nil, nil, err
	}
	// defer db.Close()

	processedBlockRepository, err := infra.NewPgProcessedBlock(db)
	if err != nil {
		return nil, nil, err
	}

	processorService, err := service.NewProcessorService(kafkaClient, codec, eventFilter, transactionMessageFilter)
	if err != nil {
		return nil, nil, err
	}

	eventsProducer, err := usecase.NewEventsProducer(processorService, alloraClient, processedBlockRepository, 0,
		cfg.Producer.BlockRefreshInterval, cfg.Producer.RateLimitInterval, cfg.Producer.NumWorkers)
	if err != nil {
		return nil, nil, err
	}

	transactionsProducer, err := usecase.NewTransactionsProducer(processorService, alloraClient, processedBlockRepository, 0,
		cfg.Producer.BlockRefreshInterval, cfg.Producer.RateLimitInterval, cfg.Producer.NumWorkers)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if db != nil {
			db.Close()
			log.Info().Msg("Database connection closed")
		}
		// Add other cleanup tasks if necessary
	}

	return app.NewApp(eventsProducer, transactionsProducer), cleanup, nil
}
