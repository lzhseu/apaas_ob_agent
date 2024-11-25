package logger

import (
	"testing"
)

func TestConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger()
	logger.Debug("this is a debug log", "foo", "bar")
	logger.Info("this is a info log", "foo", "bar")
	logger.Warn("this is a warn log", "foo", "bar")
	logger.Error("this is a error log", "foo", "bar")
}

func TestFileLogger(t *testing.T) {
	logger := NewFileLogger("test.log", HourDur, 0, 1)
	logger.Debug("this is a debug log", "foo", "bar")
	logger.Info("this is a info log", "foo", "bar")
	logger.Warn("this is a warn log", "foo", "bar")
	logger.Error("this is a error log", "foo", "bar")
}

func TestLokiLogger(t *testing.T) {
	logger := NewLokiLogger("http://10.37.107.200:3100", map[string]string{"service_name": "apaas_agent"})
	logger.Warn("loki push. this is a warn log 33", "foo", "bar2")
}
