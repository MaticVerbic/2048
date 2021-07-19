package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Log *logrus.Entry
}

func New() *Config {
	return &Config{
		Log: initLog(),
	}
}

func initLog() *logrus.Entry {
	level, err := logrus.ParseLevel("trace")
	if err != nil {
		logrus.Fatal("invalid logrus level")
	}
	logrus.SetLevel(level)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
		ForceColors:     true,
	})

	logrus.WithField("logrus_level", level).Trace("log init")
	return logrus.NewEntry(logrus.New())
}
