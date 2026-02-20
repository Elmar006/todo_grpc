package logger

import "github.com/sirupsen/logrus"

var Log = logrus.New()

func init() {

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func L() *logrus.Logger {
	return Log
}
