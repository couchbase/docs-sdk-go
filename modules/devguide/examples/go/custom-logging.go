package main

import (
	"os"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/sirupsen/logrus"
)

// #tag::loggerwrapper[]
type MyLogrusLogger struct {
	logger *logrus.Logger
}

// The logrus Log function doesn't match the gocb Log function so we need to do a bit of marshalling.
func (logger *MyLogrusLogger) Log(level gocb.LogLevel, offset int, format string, v ...interface{}) error {
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
	logger.logger.Logf(logrusLevel, format, v...)
	return nil
}

// #end::loggerwrapper[]

func main() {
	// #tag::creation[]
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	gocb.SetLogger(&MyLogrusLogger{
		logger: logger,
	})
	// #end::creation[]

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}

	cluster, err := gocb.Connect("couchbase://localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("default")
	col := bucket.DefaultCollection()

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	_, err = col.Upsert("mylogger", "logs", nil)
	if err != nil {
		panic(err)
	}

	err = cluster.Close(nil)
	if err != nil {
		panic(err)
	}
}
