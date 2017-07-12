package log

import (
	"github.com/Sirupsen/logrus"
	"path"
	"runtime"
	"strings"
)

var (
	logger *logrus.Logger
)

type Config struct {
	Level string
}

type Fields map[string]interface{}

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

func contextFields(lvl ...int) Fields {
	level := 2
	if len(lvl) == 1 {
		level = lvl[0]
	}
	pc, file, line, _ := runtime.Caller(level)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""

	if len(parts) >= 0 && pl-2 < len(parts) {
		if parts[pl-2][0] == '(' {
			packageName = strings.Join(parts[0:pl-2], ".")
		} else {
			packageName = strings.Join(parts[0:pl-1], ".")
		}

		pkgs := strings.Split(packageName, "/shuttletracker/")
		if len(pkgs) > 1 {
			packageName = pkgs[1]
		}
	}

	return Fields{
		"package": packageName,
		"file":    fileName,
		"line":    line,
	}
}

func WithField(f string, v interface{}) *logrus.Entry {
	return logger.WithField(f, v)
}

func WithFields(f ...Fields) *logrus.Entry {
	if len(f) == 0 {
		return logger.WithFields(logrus.Fields{})
	}
	e := logrus.NewEntry(logger)
	for i := 0; i < len(f); i++ {
		e = e.WithFields(logrus.Fields(f[i]))
	}
	return e
}

func WithError(err error) *logrus.Entry {
	return WithFields(contextFields()).WithField("error", err)
}

func Error(args ...interface{}) {
	WithFields(contextFields()).Error(args...)
}

func Errorf(format string, args ...interface{}) {
	WithFields(contextFields()).Errorf(format, args...)
}

func Warn(args ...interface{}) {
	WithFields(contextFields()).Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	WithFields(contextFields()).Warnf(format, args...)
}

func Info(args ...interface{}) {
	WithFields(contextFields()).Info(args...)
}

func Infof(format string, args ...interface{}) {
	WithFields(contextFields()).Infof(format, args...)
}

func Debug(args ...interface{}) {
	WithFields(contextFields()).Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	WithFields(contextFields()).Debugf(format, args...)
}
