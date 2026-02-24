package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "20060102 15:04:05",
	})
	Log.SetLevel(logrus.InfoLevel)
}

func L() *logrus.Logger {
	return Log
}
