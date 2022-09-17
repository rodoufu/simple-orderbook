package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	formatter := &logrus.TextFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true

	logger := logrus.New()
	logger.SetFormatter(formatter)
	logger.SetOutput(os.Stderr)
	log := logger.WithFields(logrus.Fields{})

	log.Info("staring service")
}
