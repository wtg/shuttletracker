package log

import (
	"path"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
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

// NewConfig creates a new Config.
func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		Level: "info",
	}
	v.SetDefault("log.level", cfg.Level)
	return cfg
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
	pkgs := strings.Split(runtime.FuncForPC(pc).Name(), "/shuttletracker/")
	packageName := strings.Split(pkgs[1], ".")[0]

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
