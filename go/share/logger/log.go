package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	GlobalLogger *zap.Logger
	Sugar        *zap.SugaredLogger
)

func NewLogger() (*zap.Logger, error) {
	// 创建配置
	config := zap.NewProductionConfig()

	// 自定义配置
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.NameKey = "logger"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 设置日志级别
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	return config.Build()
}

func Init(mode string) error {
	var (
		logger *zap.Logger
		err    error
	)

	if mode == "dev" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = NewLogger()
	}
	if err != nil {
		return err
	}

	GlobalLogger = logger
	Sugar = logger.Sugar()
	return nil
}

func Sync() {
	if GlobalLogger != nil {
		_ = GlobalLogger.Sync()
	}
}
