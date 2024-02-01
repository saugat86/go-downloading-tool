package util

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func init() {
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
}
