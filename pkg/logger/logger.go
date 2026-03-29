package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type Config struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

func New(cfg *Config) (*Logger, error) {
	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse level: %v", err)
	}
	logger.SetLevel(level)

	if cfg.Path != "" {
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("file doesn't exist: %v", err)
		}
		logger.SetOutput(file)
	}

	return &Logger{logger}, nil
}

func NewTestLogger() *Logger {
	log := logrus.New()
	log.SetOutput(io.Discard)
	return &Logger{log}
}
