package main

// #tag::connect[]
import (
	"fmt"

	gocb "github.com/couchbase/gocb/v2"
)

// #tag::connect[]
func main() {
	cluster, err := gocb.Connect(
		"localhost",
		&gocb.ClusterOptions{
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
	// #end::bucket[]

	// #tag::named-bucket[]
	// get a bucket reference
	bucket := cluster.Bucket("travel-sample")
	// #end::named-bucket[]

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
	fmt.Println(upsertResult)

	// Get Document
	getResult, err := collection.Get("my-document", &gocb.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(getResult)
	// #end::upsert-get[]
}
