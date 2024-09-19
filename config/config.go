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
	Monitor           MonitorConfig            `mapstructure:"monitor" validate:"required"`
}

type KafkaConfig struct {
	Seeds    []string `mapstructure:"seeds" validate:"required,gt=0,dive,hostname_port"` // Kafka seeds
	User     string   `mapstructure:"user" validate:"required,min=1"`                    // Kafka user
	Password string   `mapstructure:"password" validate:"required,min=1"`                // Kafka password
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
	Types []string `mapstructure:"types" validate:"required,gt=0,dive,min=1"` // Types of transactions to filter
}

type LogConfig struct {
	Level int8 `mapstructure:"level" validate:"required"` // Log level
}

type MonitorConfig struct {
	BlockRefreshInterval time.Duration `mapstructure:"block_refresh_interval" validate:"required"` // Block refresh interval
	RateLimitInterval    time.Duration `mapstructure:"rate_limit_interval" validate:"required"`    // Rate limit interval
}
