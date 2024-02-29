package exceptions

import (
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

type logger struct {
	*logrus.Logger
}

var Log logger

func init() {
	Log = logger{logrus.New()}
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
		FullTimestamp:   true,
	})
	Log.SetLevel(logrus.DebugLevel)
	Log.SetOutput(os.Stdout)
}

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	Log.Printf("%v: %v", msg, time.Since(start))
}

func (l logger) WrapError(err error) *logrus.Entry {
	var file string
	var line int
	var caller string

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok || detail != nil {
		caller = detail.Name()
	}
	pe := PkgErrorEntry{Entry: logrus.WithFields(logrus.Fields{
		"Caller": caller,
		"File":   file,
		"Line":   line,
	})}
	return pe.WithError(err)
}
