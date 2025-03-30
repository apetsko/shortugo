package logging

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestLogger_Methods(t *testing.T) {
	logger, _ := New(zapcore.DebugLevel)

	testCases := []struct {
		name    string
		logFunc func()
	}{
		{"Debug", func() { logger.Debug("debug message", "key", "value") }},
		{"Info", func() { logger.Info("info message", "key", "value") }},
		{"Error", func() { logger.Error("error message", "key", "value") }},
		{"Printf", func() { logger.Printf("formatted %s", "message") }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.logFunc()
		})
	}
}
