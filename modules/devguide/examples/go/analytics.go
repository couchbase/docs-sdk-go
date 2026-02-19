package main

import (
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	// #tag::query[]
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	// For Server versions 6.5 or later you do not need to open a bucket here
	b := cluster.Bucket("travel-sample")

	// We wait until the bucket is definitely connected and setup.
	// For Server versions 6.5 or later if we hadn't opened a bucket then we could use cluster.WaitUntilReady here.
	err = b.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	results, err := cluster.AnalyticsQuery("SELECT \"hello\" as greeting;", nil)
	if err != nil {
		panic(err)
	}

	var greeting interface{}
	for results.Next() {
		err := results.Row(&greeting)
		if err != nil {
			panic(err)
		}
		fmt.Println(greeting)
	}

	// always check for errors after iterating.
	err = results.Err()
	if err != nil {
		panic(err)
	}
	// #end::query[]

	// #tag::options[]
	results, err = cluster.AnalyticsQuery(
		"select airportname, country from airports where country = 'France';",
		&gocb.AnalyticsOptions{
			Priority: true,
			Timeout:  100 * time.Second,
		},
	)
	// #end::options[]

	cluster.Close(nil)
}
