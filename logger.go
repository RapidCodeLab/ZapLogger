package zaplogger

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func New() *ZapLogger {
	return &ZapLogger{}
}

func (l *ZapLogger) Infof(template string, args ...any) {
	l.logger.Sugar().Infof(template, args)
}

func (l *ZapLogger) Debugf(template string, args ...any) {
	l.logger.Sugar().Debugf(template, args)
}

func (l *ZapLogger) Warnf(template string, args ...any) {
	l.logger.Sugar().Warnf(template, args)
}

func (l *ZapLogger) Errorf(template string, args ...any) {
	l.logger.Sugar().Errorf(template, args)
}

func (l *ZapLogger) Fatalf(template string, args ...any) {
	l.logger.Sugar().Fatalf(template, args)
}

func (l *ZapLogger) Panicf(template string, args ...any) {
	l.logger.Sugar().Panicf(template, args)
}
