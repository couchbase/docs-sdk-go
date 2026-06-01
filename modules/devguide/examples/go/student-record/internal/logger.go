package internal

import (
	"os"

	"github.com/couchbase/gocb/v2"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	wrapped *logrus.Logger
}

func NewLogger(level logrus.Level) *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(level)

	return &Logger{
		wrapped: logger,
	}
}

// The logrus Log function doesn't match the gocb Log function so we need to do a bit of marshalling.
func (logger *Logger) Log(level gocb.LogLevel, offset int, format string, v ...interface{}) error {
	// We need to do some conversion between gocb and logrus levels as they don't match up.
	var logrusLevel logrus.Level
	switch level {
	case gocb.LogError:
		logrusLevel = logrus.ErrorLevel
	case gocb.LogWarn:
		logrusLevel = logrus.WarnLevel
	case gocb.LogInfo:
		logrusLevel = logrus.InfoLevel
	case gocb.LogDebug:
		logrusLevel = logrus.DebugLevel
	case gocb.LogTrace:
		logrusLevel = logrus.TraceLevel
	case gocb.LogSched:
		logrusLevel = logrus.TraceLevel
	case gocb.LogMaxVerbosity:
		logrusLevel = logrus.TraceLevel
	}

	// Send the data to the logrus Logf function to make sure that it gets formatted correctly.
	logger.wrapped.Logf(logrusLevel, format, v...)
	return nil
}
