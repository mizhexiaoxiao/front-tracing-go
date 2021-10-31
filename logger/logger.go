package logger

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logTmFmtWithMS = "2006-01-02 15:04:05.000"
)

var (
	logger        *zap.Logger
	infoLogger    *log.Logger
	debugLogger   *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func init() {
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(logTmFmtWithMS) + "]")
	}
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     customTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"}, //[]string{"/tmp/zap.log"},
		ErrorOutputPaths: []string{"stdout"},
		// InitialFields: map[string]interface{}{
		// 	"app": "test",
		// },
	}
	logger, _ = cfg.Build()
	infoLogger, _ = zap.NewStdLogAt(logger, zapcore.InfoLevel)
	debugLogger, _ = zap.NewStdLogAt(logger, zapcore.DebugLevel)
	warningLogger, _ = zap.NewStdLogAt(logger, zapcore.WarnLevel)
	errorLogger, _ = zap.NewStdLogAt(logger, zapcore.ErrorLevel)
}

func InfoLogger() *log.Logger {
	return infoLogger
}

func DebugLogger() *log.Logger {
	return debugLogger
}

func WarningLogger() *log.Logger {
	return warningLogger
}

func ErrorLogger() *log.Logger {
	return errorLogger
}
