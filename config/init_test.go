package config_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/allora-network/allora-producer/config"
	"github.com/allora-network/allora-producer/config/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInitConfig_Success(t *testing.T) {
	mockLoader := new(mocks.Loader)

	// Define expectations
	mockLoader.On("AddConfigPath", "./").Once()
	mockLoader.On("AddConfigPath", "./config").Once()
	mockLoader.On("AddConfigPath", "../config").Once()
	mockLoader.On("AddConfigPath", "../../config").Once()
	mockLoader.On("SetConfigName", "config").Once()
	mockLoader.On("SetConfigType", "yaml").Once()
	mockLoader.On("SetEnvPrefix", "ALLORA").Once()
	mockLoader.On("SetEnvKeyReplacer", mock.AnythingOfType("*strings.Replacer")).Once()
	mockLoader.On("AutomaticEnv").Once()
	mockLoader.On("ReadInConfig").Return(nil).Once()
	mockLoader.On("ConfigFileUsed").Return("/path/to/config.yaml").Once()

	// Mock Unmarshal
	expectedConfig := config.Config{
		Kafka: config.KafkaConfig{
			Seeds:    []string{"seed1", "seed2"},
			User:     "user",
			Password: "password",
		},
		Allora: config.AlloraConfig{
			RPC:     "http://allora-rpc.com",
			Timeout: 5000,
		},
	}
	mockLoader.On("Unmarshal", mock.AnythingOfType("*config.Config")).Run(func(args mock.Arguments) {
		arg, ok := args.Get(0).(*config.Config)
		require.True(t, ok)
		*arg = expectedConfig
	}).Return(nil).Once()

	// Initialize config provider with mock
	cfgProvider := config.NewProvider(mockLoader)

	// Call InitConfig
	cfg, err := cfgProvider.InitConfig()

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, cfg)

	// Assert expectations
	mockLoader.AssertExpectations(t)
}

func TestInitConfig_ReadConfigError(t *testing.T) {
	mockLoader := new(mocks.Loader)

	// Define expectations
	mockLoader.On("AddConfigPath", "./").Once()
	mockLoader.On("AddConfigPath", "./config").Once()
	mockLoader.On("AddConfigPath", "../config").Once()
	mockLoader.On("AddConfigPath", "../../config").Once()
	mockLoader.On("SetConfigName", "config").Once()
	mockLoader.On("SetConfigType", "yaml").Once()
	mockLoader.On("SetEnvPrefix", "ALLORA").Once()
	mockLoader.On("SetEnvKeyReplacer", mock.AnythingOfType("*strings.Replacer")).Once()
	mockLoader.On("AutomaticEnv").Once()
	mockLoader.On("ReadInConfig").Return(errors.New("config file not found")).Once()

	// Initialize config provider with mock
	cfgProvider := config.NewProvider(mockLoader)

	// Call InitConfig
	cfg, err := cfgProvider.InitConfig()

	// Assertions
	require.Error(t, err)
	assert.Equal(t, config.Config{}, cfg)
	assert.Contains(t, err.Error(), "error reading config file")

	// Assert expectations
	mockLoader.AssertExpectations(t)
}

func TestInitConfig_UnmarshalError(t *testing.T) {
	mockLoader := new(mocks.Loader)

	// Define expectations
	mockLoader.On("AddConfigPath", "./").Once()
	mockLoader.On("AddConfigPath", "./config").Once()
	mockLoader.On("AddConfigPath", "../config").Once()
	mockLoader.On("AddConfigPath", "../../config").Once()
	mockLoader.On("SetConfigName", "config").Once()
	mockLoader.On("SetConfigType", "yaml").Once()
	mockLoader.On("SetEnvPrefix", "ALLORA").Once()
	mockLoader.On("SetEnvKeyReplacer", mock.AnythingOfType("*strings.Replacer")).Once()
	mockLoader.On("AutomaticEnv").Once()
	mockLoader.On("ReadInConfig").Return(nil).Once()
	mockLoader.On("ConfigFileUsed").Return("/path/to/config.yaml").Once()
	mockLoader.On("Unmarshal", mock.AnythingOfType("*config.Config")).Return(errors.New("error unmarshaling config")).Once()

	// Initialize config provider with mock
	cfgProvider := config.NewProvider(mockLoader)

	// Call InitConfig
	cfg, err := cfgProvider.InitConfig()

	// Assertions
	require.Error(t, err)
	assert.Equal(t, config.Config{}, cfg)
	assert.Contains(t, err.Error(), "error unmarshaling config")

	// Assert expectations
	mockLoader.AssertExpectations(t)
}

func TestValidate(t *testing.T) {
	cfg := config.Config{
		Database: config.DatabaseConfig{
			URL: "postgres://user:password@localhost:5432/dbname",
		},
		Kafka: config.KafkaConfig{
			Seeds:    []string{"localhost:8000", "localhost:8001"},
			User:     "user",
			Password: "password",
		},
		Allora: config.AlloraConfig{
			RPC:     "http://allora-rpc.com",
			Timeout: 5000,
		},
		KafkaTopicRouter: []config.KafkaTopicRouterConfig{
			{
				Name:  "topic1",
				Types: []string{"event1", "transaction1"},
			},
		},
		FilterEvent: config.FilterEventConfig{
			Types: []string{"event1", "event2"},
		},
		FilterTransaction: config.FilterTransactionConfig{
			Types: []string{"transaction1", "transaction2"},
		},
		Log: config.LogConfig{
			Level: 1,
		},
		Producer: config.ProducerConfig{
			BlockRefreshInterval: 1000,
			RateLimitInterval:    1000,
			NumWorkers:           10,
		},
	}

	err := config.ValidateConfig(&cfg)
	require.NoError(t, err)
}

func TestInitConfig_Load(t *testing.T) {
	// If config.yaml doesn't exist, copy config.yaml.example to config.yaml
	configCreated := false
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		input, err := os.ReadFile("config.example.yaml")
		require.NoError(t, err)

		err = os.WriteFile("config.yaml", input, 0600)
		require.NoError(t, err)
		configCreated = true
	}

	cfgProvider := config.NewProvider(config.NewViperAdapter())
	cfg, err := cfgProvider.InitConfig()
	require.NoError(t, err)
	assert.NotEmpty(t, cfg)

	// Validate the config
	err = config.ValidateConfig(&cfg)
	require.NoError(t, err)

	// Remove the created config.yaml file if it was created in this test
	if configCreated {
		err = os.Remove("config.yaml")
		require.NoError(t, err)
	}
}

func TestViperAdapter(t *testing.T) {
	adapter := config.NewViperAdapter()
	adapter.AddConfigPath("./")
	adapter.SetConfigName("config.example")
	adapter.SetConfigType("yaml")
	adapter.SetEnvPrefix("ALLORA")
	adapter.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	adapter.AutomaticEnv()

	err := adapter.ReadInConfig()
	require.NoError(t, err)

	configFileUsed := adapter.ConfigFileUsed()
	wd, err := os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, wd+"/config.example.yaml", configFileUsed)

	cfg := config.Config{}
	err = adapter.Unmarshal(&cfg)
	require.NoError(t, err)

	assert.NotEmpty(t, cfg)
}
