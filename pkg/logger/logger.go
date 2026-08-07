package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.Out = os.Stdout

	// Set format
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set level
	if os.Getenv("APP_ENV") == "production" {
		Log.SetLevel(logrus.InfoLevel)
	} else {
		Log.SetLevel(logrus.DebugLevel)
	}
}
