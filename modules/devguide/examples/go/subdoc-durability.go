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
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("default")
	collection := bucket.DefaultCollection()

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	// Observe based
	// #tag::traddurability[]
	mops := []gocb.MutateInSpec{
		gocb.InsertSpec("name", "mike", nil),
	}
	observeResult, err := collection.MutateIn("key", mops, &gocb.MutateInOptions{
		PersistTo:     1,
		ReplicateTo:   1,
		StoreSemantic: gocb.StoreSemanticsUpsert,
	})
	if err != nil {
		panic(err)
	}
	// #end::traddurability[]
	fmt.Println(observeResult.Cas())

	// Enhanced
	// #tag::newdurability[]
	mops = []gocb.MutateInSpec{
		gocb.InsertSpec("name", "mike", nil),
	}
	durableResult, err := collection.MutateIn("key2", mops, &gocb.MutateInOptions{
		DurabilityLevel: gocb.DurabilityLevelMajority,
		StoreSemantic:   gocb.StoreSemanticsUpsert,
	})
	if err != nil {
		panic(err)
	}
	// #end::newdurability[]
	fmt.Println(durableResult.Cas())

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
