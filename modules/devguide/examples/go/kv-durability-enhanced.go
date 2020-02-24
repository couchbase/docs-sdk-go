package main

import (
	"fmt"

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
	cluster, err := gocb.Connect("10.112.193.101", opts)
	if err != nil {
		panic(err)
	}

	// Open Bucket and collection
	bucket := cluster.Bucket("default")
	collection := bucket.DefaultCollection()

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
