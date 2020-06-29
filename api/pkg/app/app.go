package app

import (
	"strings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

)

type Base interface {
	Environment() EnvMode
	Logger() *zap.SugaredLogger
	Cleanup()
}

type Config interface {
	Base
}

type EnvMode string

const (
	Production  EnvMode = "production"
	Development EnvMode = "development"
	Test        EnvMode = "test"
)

type BaseConfig struct {
	mode   EnvMode
	logger *zap.SugaredLogger
}


func (bc *BaseConfig) Environment() EnvMode {
	return bc.mode
}

func (bc *BaseConfig) Logger() *zap.SugaredLogger {
	return bc.logger
}

func (bc *BaseConfig) Cleanup() {
	bc.logger.Sync()
}

type ApiConfig struct {
	*BaseConfig
}

func BaseConfigFromEnv() (*BaseConfig, error) {
	mode := Environment()
	var err error

	var log *zap.SugaredLogger
	if log, err = initLogger(mode); err != nil {
		return nil, err
	}

	log.With("name", "app").Infof("in %q mode ", mode)

	 bc := &BaseConfig{mode: mode, logger: log}
	return bc, nil
}

func FromEnv() (*ApiConfig, error) {
	bc, err := BaseConfigFromEnv()
	if err != nil {
		return nil, err
	}
	ApiConfig := &ApiConfig{BaseConfig: bc}

	return ApiConfig, nil
}


func Environment() EnvMode {
	mode := "development"

	switch strings.ToLower(mode) {
	case "development":
		return Development
	case "test":
		return Test
	default:
		return Production
	}
}

func initLogger(mode EnvMode) (*zap.SugaredLogger, error) {

	var log *zap.Logger
	var err error

	switch mode {
	case Production:
		log, err = zap.NewProduction()

	default:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		log, err = config.Build()
	}

	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}