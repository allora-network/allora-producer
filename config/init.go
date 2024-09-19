package config

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitConfig() (Config, error) {
	configPaths := []string{"config", "../../config", "./"}
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Automatically read configuration from environment variables
	viper.SetEnvPrefix("ALLORA") // Use a prefix for environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	log.Info().Msgf("Configuration loaded from: %s", viper.ConfigFileUsed())

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config: %w", err)
	}

	log.Info().Msg("Configuration loaded successfully")
	return cfg, nil
}
