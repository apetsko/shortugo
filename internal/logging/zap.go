package logging

import (
	"sync"

	"go.uber.org/zap"
)

type ZapLogger struct {
	log *zap.SugaredLogger
}

var instance *ZapLogger
var once sync.Once

func NewZapLogger() (*ZapLogger, error) {
	var err error

	once.Do(func() {
		logger, e := zap.NewDevelopment()
		if e != nil {
			err = e
			return
		}
		instance = &ZapLogger{
			log: logger.Sugar(),
		}
	})

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
