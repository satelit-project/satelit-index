package logging

import (
	"satelit-project/satelit-index/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const prodLogFile = "/var/log/satelit/import.log"

var defaultLogger *zap.Logger

func init() {
	logger, err := NewLogger(config.CurrentEnvironment())
	if err != nil {
		panic(err)
	}

	minLevel := int8(zapcore.InfoLevel)
	maxLevel := int8(zapcore.FatalLevel)
	for i := minLevel; i <= maxLevel; i++ {
		_, err = zap.RedirectStdLogAt(logger, zapcore.Level(i))
		if err != nil {
			panic(err)
		}
	}

	defaultLogger = logger
}

func DefaultLogger() *zap.SugaredLogger {
	return defaultLogger.Sugar()
}

func NewLogger(env config.Environment) (*zap.Logger, error) {
	var cfg zap.Config

	if env == config.ProductionEnvironment {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = append(cfg.OutputPaths, prodLogFile)
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	return cfg.Build()
}
