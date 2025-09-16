package core

import (
	"github.com/tommzn/go-log"
)

type CustomLogger struct {
	log log.Logger
}

func NewCustomLogger(log log.Logger) *CustomLogger {
	return &CustomLogger{log: log}
}

func (l *CustomLogger) Error(message string, v ...interface{}) {
	l.log.Errorf(message, v...)
}

func (l *CustomLogger) Info(message string, v ...interface{}) {
	l.log.Infof(message, v...)
}

func (l *CustomLogger) Debug(message string, v ...interface{}) {
	l.log.Debugf(message, v...)
}
