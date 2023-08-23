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
	cluster, err := gocb.Connect("your-ip", opts)
	if err != nil {
		panic(err)
	}

	// For Server versions 6.5 or later you do not need to open a bucket here
	bucket := cluster.Bucket("travel-sample")

	// We wait until the bucket is definitely connected and setup.
	// For Server versions 6.5 or later if we hadn't opened a bucket then we could use cluster.WaitUntilReady here.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	{
		fmt.Println("handle-collection")

		// #tag::handle-collection[]
		results, err := cluster.AnalyticsQuery("select airportname, country from `travel-sample`.inventory.airport where country = 'France' limit 3;", nil)
		if err != nil {
			panic(err)
		}
		// #end::handle-collection[]

		var row interface{}
		for results.Next() {
			err := results.Row(&row)
			if err != nil {
				panic(err)
			}
			fmt.Println(row)
		}
	}

	{
		fmt.Println("handle-scope")

		// tag::handle-scope[]
		scope := bucket.Scope("inventory")
		results, err := scope.AnalyticsQuery("SELECT airportname, country FROM `airport` WHERE country='France' LIMIT 2", nil)
		if err != nil {
			panic(err)
		}
		// end::handle-scope[]

		var row interface{}
		for results.Next() {
			err := results.Row(&row)
			if err != nil {
				panic(err)
			}
			fmt.Println(row)
		}
	}

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
