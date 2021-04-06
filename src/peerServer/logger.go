package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLog() {
	logrus.SetLevel(getLogLevel())
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func getLogLevel() logrus.Level {
	level := os.Getenv("LOG_LEVEL")

	switch level {
	case "debug":
		return logrus.DebugLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel
	}
}
