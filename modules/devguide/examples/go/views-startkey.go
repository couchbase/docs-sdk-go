package main

import (
	"fmt"

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

	// get a bucket reference
	bucket := cluster.Bucket("travel-sample")

	// #tag::beerview[]
	viewResult, err := bucket.ViewQuery("beer", "by_name", &gocb.ViewOptions{
		StartKey: "A",
		Limit:    10,
	})
	if err != nil {
		panic(err)
	}
	// #end::beerview[]

	for viewResult.Next() {
		row := viewResult.Row()
		fmt.Printf("Document ID: %s\n", row.ID)
		var key string
		err = row.Key(&key)
		if err != nil {
			panic(err)
		}

		var beer interface{}
		err = row.Value(&beer)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Beer named %s has value %v\n", key, beer)
	}

	// always check for errors after iterating
	err = viewResult.Err()
	if err != nil {
		panic(err)
	}

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
