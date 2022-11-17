package logger

import "github.com/sirupsen/logrus"

type LogrusLoggerWrapper struct {
	logger *logrus.Logger
}

func NewLogrusLogger(logger *logrus.Logger) *LogrusLoggerWrapper {
	return &LogrusLoggerWrapper{logger: logger}
}

func (l *LogrusLoggerWrapper) Flush() {
}

func (l *LogrusLoggerWrapper) Tracef(format string, params ...interface{}) {
	l.logger.Tracef(format, params)
}

func (l *LogrusLoggerWrapper) Infof(format string, params ...interface{}) {
	l.logger.Infof(format, params)
}

func (l *LogrusLoggerWrapper) Debugf(format string, params ...interface{}) {
	l.logger.Debugf(format, params)
}

func (l *LogrusLoggerWrapper) Warnf(format string, params ...interface{}) {
	l.logger.Warnf(format, params)
}

func (l *LogrusLoggerWrapper) Errorf(format string, params ...interface{}) {
	l.logger.Errorf(format, params)
}
