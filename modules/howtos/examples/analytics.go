package main

import (
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.Connect("10.112.194.101", opts)
	if err != nil {
		panic(err)
	}

	// #tag::query[]
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
