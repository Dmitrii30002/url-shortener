package config

import (
	"os"

	"github.com/Dmitrii30002/url-shortener/pkg/logger"
	"github.com/Dmitrii30002/url-shortener/pkg/storage/postgres"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	StorageType string `yaml:"storage_type"`

	BaseURL string `yaml:"base_url"`

	PostgresCfg postgres.Config

	LoggerCfg logger.Config `yaml:"logger"`
}

func GetConfig(path string) (*Config, error) {
	config := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	config.PostgresCfg.Host = os.Getenv("POSTGRES_HOST")
	config.PostgresCfg.Port = os.Getenv("POSTGRES_PORT")
	config.PostgresCfg.Name = os.Getenv("POSTGRES_DB")
	config.PostgresCfg.User = os.Getenv("POSTGRES_USER")
	config.PostgresCfg.Password = os.Getenv("POSTGRES_PASSWORD")
	config.PostgresCfg.SSLMode = os.Getenv("POSTGRES_SSLMODE")

	return config, nil
}

func GetTestConfig() (*Config, error) {
	config := &Config{}

	config.PostgresCfg.Host = os.Getenv("TEST_POSTGRES_HOST")
	config.PostgresCfg.Port = os.Getenv("TEST_POSTGRES_PORT")
	config.PostgresCfg.Name = os.Getenv("TEST_POSTGRES_DB")
	config.PostgresCfg.User = os.Getenv("TEST_POSTGRES_USER")
	config.PostgresCfg.Password = os.Getenv("TEST_POSTGRES_PASSWORD")
	config.PostgresCfg.SSLMode = os.Getenv("TEST_POSTGRES_SSLMODE")

	return config, nil
}
