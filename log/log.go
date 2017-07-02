package log

import (
	"github.com/Sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

type Config struct {
	Level string
}

func init() {
	logger = logrus.New()
}

func NewConfig() *Config {
	return &Config{
		Level: "info",
	}
}

func SetLevel(level string) {
	parsed, err := logrus.ParseLevel(level)
	if err != nil {
		Error(err)
		return
	}
	logger.Level = parsed
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
