package logging

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	log *zap.SugaredLogger
}

func NewZapLogger() (*ZapLogger, error) {
	var err error
	logger, e := zap.NewDevelopment()
	if e != nil {
		err = e
		return nil, err
	}
	instance := &ZapLogger{
		log: logger.Sugar(),
	}
	return instance, err
}

func (l *ZapLogger) Info(message string, keysAndValues ...interface{}) {
	l.log.Infow(message, keysAndValues...)
}

func (l *ZapLogger) Error(message string, keysAndValues ...interface{}) {
	l.log.Errorw(message, keysAndValues...)
}

func (l *ZapLogger) Fatal(message string, keysAndValues ...interface{}) {
	l.log.Fatalw(message, keysAndValues...)
}
