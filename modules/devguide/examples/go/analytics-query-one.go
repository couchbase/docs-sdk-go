package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
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

	// #tag::one[]
	rows, err := cluster.AnalyticsQuery("select airportname, country from airports where country = 'France' LIMIT 1;", nil)
	// check query was successful
	if err != nil {
		panic(err)
	}

	var result interface{}
	err = rows.One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// #end::one[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
