package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log writes a logger line
	log *zap.Logger
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	return log
}

// Info logs an info level
func Info(msg string, tags ...zap.Field) {
	log.Info(msg, tags...)
	log.Sync()
}

// Error logs an error level
func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Error(msg, tags...)
	log.Sync()
}
