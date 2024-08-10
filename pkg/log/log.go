package log

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logger *zap.Logger
}

type LoggerMetaData struct {
	LogLevel   string
	Service    string
	AppVersion string
}

type LogField struct {
	Key   string
	Value any
}

var (
	globalLogger logger
	once         sync.Once
)

func InitLogger(metaData LoggerMetaData) {
	once.Do(func() {
		zapConfig := zap.NewProductionConfig()
		zapConfig.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

		zapOptions := []zap.Option{
			zap.Fields(
				zap.String("service", metaData.Service),
				zap.String("app_version", metaData.AppVersion),
			),
		}

		stackTraceLevel := zap.ErrorLevel
		if metaData.LogLevel != "" {
			logLevel, err := zapcore.ParseLevel(metaData.LogLevel)
			if err != nil {
				panic(fmt.Errorf("failed parse log level, err detail: %v", err))
			}

			stackTraceLevel = logLevel
			zapOptions = append(zapOptions, zap.IncreaseLevel(logLevel))
		}

		zapOptions = append(zapOptions, zap.AddStacktrace(stackTraceLevel))

		globalLogger.logger = zap.Must(zapConfig.Build(zapOptions...))
	})
}

func GetLogger() logger {
	if globalLogger.logger == nil {
		return logger{
			logger: zap.NewNop(),
		}
	}

	return globalLogger
}

func Info(ctx context.Context, msg string, fields ...LogField) {
	// TODO: add fields
	GetLogger().logger.Info(msg)
}

func Error(ctx context.Context, msg string, err error, fields ...LogField) {
	// TODO: add fields
	GetLogger().logger.Error(msg, zap.Error(err))
}

func Fatal(ctx context.Context, msg string, err error, fields ...LogField) {
	// TODO: add fields
	GetLogger().logger.Fatal(msg, zap.Error(err))
}
