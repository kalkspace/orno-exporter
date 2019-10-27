package orno

import (
	"log"

	"github.com/sirupsen/logrus"
)

type logWriter struct {
	logger logrus.FieldLogger
}

func (l *logWriter) Write(msg []byte) (int, error) {
	l.logger.Debug(string(msg))
	return len(msg), nil
}

func modbusLogger(logger logrus.FieldLogger) *log.Logger {
	return log.New(&logWriter{logger}, "", 0)
}
