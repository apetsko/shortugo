package logging

import (
	"go.uber.org/zap"
)

// ZapLogger — конкретная реализация Logger, использующая zap
type ZapLogger struct {
	log *zap.SugaredLogger
}

// NewZapLogger создаёт новый экземпляр ZapLogger
func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{
		log: logger.Sugar(),
	}, nil
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
