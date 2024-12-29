package logging

import (
	"go.uber.org/zap"
)

var zapLog zap.SugaredLogger

func Start() (err error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			return
		}
	}()

	zapLog = *logger.Sugar()
	return
}

func Info(message string, keysAndValues ...interface{}) {
	zapLog.Infow(message, keysAndValues...)
}
func Fatal(message string, keysAndValues ...interface{}) {
	zapLog.Infow(message, keysAndValues...)
}
