package main

import (
	"github.com/couchbase/gocb/v2"
	"time"
)

func main() {
	// #tag::config[]
	tracerOpts := &gocb.ThresholdLoggingOptions{
		Interval:   10 * time.Second,
		SampleSize: 10,
	}
	tracer := gocb.NewThresholdLoggingTracer(tracerOpts)

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
		Tracer: tracer,
	}
	// #end::config[]
	throwaway(opts)
}
