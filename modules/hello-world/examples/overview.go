package main

import (
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	cluster, err := gocb.Connect(
		"localhost",
		gocb.ClusterOptions{
			Username: "Administrator",
			Password: "password",
		})
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("travel-sample")
	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	// tag::overview[]
	// get a default collection reference
	collection := bucket.DefaultCollection()

	// for a named collection and scope
	scope := bucket.Scope("inventory")
	collection = scope.Collection("airport")
	// end::overview[]

	_, err = collection.Upsert("airport_111", "hello-world", &gocb.UpsertOptions{})
	if err != nil {
		panic(err)
	}
}
