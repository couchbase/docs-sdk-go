package main

// #tag::connect[]
import (
	"fmt"

	gocb "github.com/couchbase/gocb/v2"
)

// #tag::connect[]
func main() {
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
	// #end::connect[]

	// #tag::bucket[]
	// get a bucket reference
	bucket := cluster.Bucket("bucket-name", &gocb.BucketOptions{})
	// #end::bucket[]

	// #tag::named-bucket[]
	// get a bucket reference
	bucket := cluster.Bucket("travel-sample", nil)
	// #end::named-bucket[]

	// #tag::collection[]
	// get a collection reference
	collection := bucket.DefaultCollection()
	// for a named collection and scope
	// collection := bucket.Scope("my-scope").Collection("my-collection", &gocb.CollectionOptions{})
	// #end::collection[]

	// #tag::upsert-get[]
	// Upsert Document
	upsertResult, err := collection.Upsert("my-document", map[string]string{"name": "mike"}, &gocb.UpsertOptions{})
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
