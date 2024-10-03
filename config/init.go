package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//go:generate mockery --all

// ConfigLoader defines the methods required for configuration loading.
// This interface abstracts the Viper library to facilitate unit testing.
type Loader interface {
	AddConfigPath(path string)
	SetConfigName(name string)
	SetConfigType(configType string)
	SetEnvPrefix(prefix string)
	SetEnvKeyReplacer(replacer *strings.Replacer)
	AutomaticEnv()
	ReadInConfig() error
	ConfigFileUsed() string
	Unmarshal(rawVal interface{}) error
}

// ViperAdapter is a concrete implementation of ConfigLoader using Viper.
type ViperAdapter struct {
	*viper.Viper
}

// NewViperAdapter creates a new instance of ViperAdapter.
func NewViperAdapter() Loader {
	return &ViperAdapter{viper.New()}
}

// AddConfigPath adds a path for Viper to search for the configuration file.
func (v *ViperAdapter) AddConfigPath(path string) {
	v.Viper.AddConfigPath(path)
}

// SetConfigName sets the name of the configuration file without the extension.
func (v *ViperAdapter) SetConfigName(name string) {
	v.Viper.SetConfigName(name)
}

// SetConfigType sets the type of the configuration file (e.g., "yaml").
func (v *ViperAdapter) SetConfigType(configType string) {
	v.Viper.SetConfigType(configType)
}

// SetEnvPrefix sets the prefix for environment variables.
func (v *ViperAdapter) SetEnvPrefix(prefix string) {
	v.Viper.SetEnvPrefix(prefix)
}

// SetEnvKeyReplacer sets the key replacer for environment variables.
func (v *ViperAdapter) SetEnvKeyReplacer(replacer *strings.Replacer) {
	v.Viper.SetEnvKeyReplacer(replacer)
}

// AutomaticEnv enables Viper to automatically read configuration from environment variables.
func (v *ViperAdapter) AutomaticEnv() {
	v.Viper.AutomaticEnv()
}

// ReadInConfig reads the configuration from the specified file.
func (v *ViperAdapter) ReadInConfig() error {
	return v.Viper.ReadInConfig()
}

// ConfigFileUsed returns the file path of the configuration file used.
func (v *ViperAdapter) ConfigFileUsed() string {
	return v.Viper.ConfigFileUsed()
}

// Unmarshal unmarshals the configuration into the provided structure.
func (v *ViperAdapter) Unmarshal(rawVal interface{}) error {
	return v.Viper.Unmarshal(rawVal)
}

type Provider interface {
	InitConfig() (Config, error)
}

type config struct {
	loader Loader
}

func NewProvider(loader Loader) Provider {
	return &config{
		loader: loader,
	}
}

func (c *config) InitConfig() (Config, error) {
	configPaths := []string{"./", "./config", "../config", "../../config"}
	for _, path := range configPaths {
		c.loader.AddConfigPath(path)
	}
	c.loader.SetConfigName("config")
	c.loader.SetConfigType("yaml")

	// Automatically read configuration from environment variables
	c.loader.SetEnvPrefix("ALLORA") // Use a prefix for environment variables
	c.loader.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.loader.AutomaticEnv()

	if err := c.loader.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	log.Info().Msgf("Configuration loaded from: %s", c.loader.ConfigFileUsed())

	var cfg Config
	if err := c.loader.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config: %w", err)
	}

	log.Info().Msg("Configuration loaded successfully")
	return cfg, nil
}

func ValidateConfig(cfg *Config) error {
	validate := validator.New()
	return validate.Struct(cfg)
}
