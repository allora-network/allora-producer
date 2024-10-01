package config

import "time"

type Config struct {
	Database          DatabaseConfig           `mapstructure:"database" validate:"required"`
	Kafka             KafkaConfig              `mapstructure:"kafka" validate:"required"`
	Allora            AlloraConfig             `mapstructure:"allora" validate:"required"`
	KafkaTopicRouter  []KafkaTopicRouterConfig `mapstructure:"kafka_topic_router" validate:"required,gt=0,dive"`
	FilterEvent       FilterEventConfig        `mapstructure:"filter_event" validate:"required"`
	FilterTransaction FilterTransactionConfig  `mapstructure:"filter_transaction" validate:"required"`
	Log               LogConfig                `mapstructure:"log" validate:"required"`
	Producer          ProducerConfig           `mapstructure:"producer" validate:"required"`
}

type KafkaConfig struct {
	Seeds         []string `mapstructure:"seeds" validate:"required,gt=0,dive,hostname_port"` // Kafka seeds
	User          string   `mapstructure:"user" validate:"required,min=1"`                    // Kafka user
	Password      string   `mapstructure:"password" validate:"required,min=1"`                // Kafka password
	NumPartitions int32    `mapstructure:"num_partitions" validate:"required,min=1"`          // Number of partitions for the Kafka topic
}

type DatabaseConfig struct {
	URL string `mapstructure:"url" validate:"required,url"` // URL of the database
}

type AlloraConfig struct {
	RPC     string `mapstructure:"rpc" validate:"required,url"`       // RPC URL of the Allora chain
	Timeout uint   `mapstructure:"timeout" validate:"required,min=1"` // Timeout in seconds for RPC requests
}

type KafkaTopicRouterConfig struct {
	Name  string   `mapstructure:"name" validate:"required"`                  // Name of the Kafka topic
	Types []string `mapstructure:"types" validate:"required,gt=0,dive,min=1"` // Types of messages to route to the topic
}

type FilterEventConfig struct {
	Types []string `mapstructure:"types" validate:"required,gt=0,dive,min=1"` // Types of events to filter
}

type FilterTransactionConfig struct {
	Types []string `mapstructure:"types" validate:"gt=0,dive,min=1"` // Types of transactions to filter
}

type LogConfig struct {
	Level int8 `mapstructure:"level" validate:"gte=-1,lte=7"` // Log level
}

type ProducerConfig struct {
	BlockRefreshInterval time.Duration `mapstructure:"block_refresh_interval" validate:"required"` // Block refresh interval
	RateLimitInterval    time.Duration `mapstructure:"rate_limit_interval" validate:"required"`    // Rate limit interval
	NumWorkers           int           `mapstructure:"num_workers" validate:"required,gt=0"`       // Number of workers to process blocks and block results
}
