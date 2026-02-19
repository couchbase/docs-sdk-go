package main

import (
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	// #tag::simple[]
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

	results, err := cluster.Query("SELECT \"hello\" as greeting;", &gocb.QueryOptions{
		// Note that we set Adhoc to true to prevent this query being run as a prepared statement.
		Adhoc: true,
	})
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

	// always check for errors after iterating
	err = results.Err()
	if err != nil {
		panic(err)
	}
	// #end::simple[]

	// tag::simple-named-scope[]
	scope := cluster.Bucket("travel-sample").Scope("inventory")
	results, err = scope.Query("SELECT x.* FROM `airline` x LIMIT 10;", &gocb.QueryOptions{})
	// check query was successful
	if err != nil {
		panic(err)
	}

	var airline interface{}
	for results.Next() {
		err := results.Row(&airline)
		if err != nil {
			panic(err)
		}
		fmt.Println(airline)
	}
	// end::simple-named-scope[]

	// always check for errors after iterating
	err = results.Err()
	if err != nil {
		panic(err)
	}
}
