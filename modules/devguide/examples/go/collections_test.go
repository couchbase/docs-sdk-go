package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

func Example_collections() {
	// Connect to Cluster
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

	bucket := cluster.Bucket("travel-sample")

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	collectionMgr := bucket.CollectionsV2()

	// create collection in default scope
	err = collectionMgr.CreateCollection("_default", "bookings", &gocb.CreateCollectionSettings{}, &gocb.CreateCollectionOptions{})
	if err != nil {
		panic(err)
	}

	// tag::collections_1[]
	bucket.Collection("bookings") // in default scope
	// end::collections_1[]

	// tag::collections_2[]
	bucket.Scope("tenant_agent_00").Collection("bookings")
	// end::collections_2[]

	fmt.Println("Done.")

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
