package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// Connect to Cluster
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	// Open Bucket and collection
	bucket := cluster.Bucket("default")
	collection := bucket.DefaultCollection()

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	type myDoc struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	document := myDoc{Foo: "bar", Bar: "foo"}
	// #tag::durability[]
	// Upsert with Durability level Majority
	durableResult, err := collection.Upsert("document-key", &document, &gocb.UpsertOptions{
		DurabilityLevel: gocb.DurabilityLevelMajority,
	})
	// #end::durability[]
	if err != nil {
		panic(err)
	}
	fmt.Println(durableResult.Cas())

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
