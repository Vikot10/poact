package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func mustCreateLogger(isDebug bool) *zap.Logger {
	devEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	prodEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	loggerConfig := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:   false,
		Encoding:      "json",
		EncoderConfig: prodEncoderConfig,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if isDebug {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		loggerConfig.Development = true
		loggerConfig.Encoding = "console"
		loggerConfig.EncoderConfig = devEncoderConfig
		loggerConfig.Sampling = nil
	}

	logger, errLogger := loggerConfig.Build()
	if errLogger != nil {
		panic(errLogger)
	}

	return logger
}
