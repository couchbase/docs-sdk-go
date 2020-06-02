package main

import (
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

// #tag::connect[]
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
	// #end::connect[]

	// #tag::bucket[]
	// get a bucket reference
	bucket := cluster.Bucket("bucket-name")

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}
	// #end::bucket[]

	// #tag::collection[]
	// get a collection reference
	collection := bucket.DefaultCollection()

	// for a named collection and scope
	// scope := bucket.Scope("my-scope")
	// collection := scope.Collection("my-collection")
	// #end::collection[]

	// #tag::upsert-get[]
	// Upsert Document
	upsertData := map[string]string{"name": "mike"}
	upsertResult, err := collection.Upsert("my-document", upsertData, &gocb.UpsertOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(upsertResult.Cas())

	// Get Document
	getResult, err := collection.Get("my-document", &gocb.GetOptions{})
	if err != nil {
		panic(err)
	}

	var myContent interface{}
	if err := getResult.Content(&myContent); err != nil {
		panic(err)
	}
	fmt.Println(myContent)
	// #end::upsert-get[]
}
