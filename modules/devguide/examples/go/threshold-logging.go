package main

import (
	"github.com/couchbase/gocb/v2"
	"time"
)

func main() {
	// tag::configure[]
	gocb.SetLogger(gocb.VerboseStdioLogger())
	tLogger := gocb.NewThresholdLoggingTracer(&gocb.ThresholdLoggingOptions{
		KVThreshold: 500 * time.Millisecond,
		SampleSize:  10,
	})
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
		Tracer: tLogger,
	}
	// end::configure[]
	connString := "localhost"
	cluster, err := gocb.Connect(connString, opts)
	if err != nil {
		panic(err)
	}

	bucketName := "travel-sample"
	// For Server versions 6.5 or later you do not need to open a bucket here
	bucket := cluster.Bucket(bucketName)

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	cluster.Close(nil)
}
