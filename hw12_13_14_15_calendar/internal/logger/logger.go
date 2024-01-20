package logger

import (
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logg interface {
	Error(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
}

func New(config config.LoggerConf) *zap.SugaredLogger {
	lavel, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return nil
	}
	logConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(lavel),
		DisableCaller:     true,
		Development:       true,
		DisableStacktrace: true,
		Encoding:          config.LogEncoding,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
	}
	logger := zap.Must(logConfig.Build()).Sugar()

	return logger
}
