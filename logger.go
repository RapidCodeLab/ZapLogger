package zaplogger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func New(c *Config) (*ZapLogger, error) {

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	generalPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.DebugLevel
	})

	cores := []zapcore.Core{}

	serviceIDField := zap.Field{
		Key:    "service",
		Type:   zapcore.StringType,
		String: c.ServiceID,
	}

	if c.StdOutLoggerEnabled {
		core := buildStdOutCore(generalPriority)
		core = core.With([]zap.Field{serviceIDField})
		cores = append(cores, core)
	}
	if c.FileLoggerEnabled {
		core := buildFileCore(highPriority, c)
		core = core.With([]zap.Field{serviceIDField})
		cores = append(cores, core)
	}

	if c.StreamLoggerEnabled {
		core := buildStreamCore(generalPriority, c)
		core = core.With([]zap.Field{serviceIDField})
		cores = append(cores, core)
	}

	if len(cores) < 1 {
		return nil, errors.New("no one logging core enabled")
	}

	tee := zapcore.NewTee(
		cores...,
	)

	return &ZapLogger{
		logger: zap.New(tee),
	}, nil
}

func (l *ZapLogger) Infof(template string, args ...any) {
	l.logger.Sugar().Infof(template, args...)
}

func (l *ZapLogger) Debugf(template string, args ...any) {
	l.logger.Sugar().Debugf(template, args...)
}

func (l *ZapLogger) Warnf(template string, args ...any) {
	l.logger.Sugar().Warnf(template, args...)
}

func (l *ZapLogger) Errorf(template string, args ...any) {
	l.logger.Sugar().Errorf(template, args...)
}

func (l *ZapLogger) Fatalf(template string, args ...any) {
	l.logger.Sugar().Fatalf(template, args...)
}

func (l *ZapLogger) Panicf(template string, args ...any) {
	l.logger.Sugar().Panicf(template, args...)
}
