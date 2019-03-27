package main

// #tag::connect[]
import (
	"fmt"

	"github.com/couchbase/gocb"
)

func main() {

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.NewCluster("localhost", opts)
	if err != nil {
		// handle err
	}
	// #end::connect[]

	// #tag::bucket[]
	// get a bucket reference
	bucket := cluster.Bucket("bucket-name", &gocb.BucketOptions{})
	// #end::bucket[]

	// #tag::collection[]
	// get a collection reference
	collection := bucket.DefaultCollection(&gocb.CollectionOptions{})
	// for a named collection and scope
	// collection := bucket.Scope("my-scope").Collection("my-collection", &gocb.CollectionOptions{})
	// #end::collection[]

	// #tag::upsert-get[]
	// Upsert Document
	upsertResult, _ := collection.Upsert("my-document", map[string]string{"name": "mike"}, &gocb.UpsertOptions{})
	fmt.Println(upsertResult)

	// Get Document
	getResult, _ := collection.Get("my-document", &gocb.GetOptions{})
	fmt.Println(getResult)
	// #end::upsert-get[]
}
