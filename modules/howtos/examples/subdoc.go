package main

import (
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	cluster, err := gocb.Connect("172.23.111.3", opts)
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

	// #tag::concurrent[]
	mops := []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("purchases.complete", 99, &gocb.ArrayAppendSpecOptions{}),
	}
	firstConcurrentResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}

	mops = []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("purchases.abandoned", 101, &gocb.ArrayAppendSpecOptions{}),
	}
	secondConcurrentResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// #end::concurrent[]
	fmt.Println(firstConcurrentResult.Cas())
	fmt.Println(secondConcurrentResult.Cas())

	_, err = collection.Upsert("player432", map[string]int{"gold": 1000}, &gocb.UpsertOptions{})
	if err != nil {
		panic(err)
	}

	// #tag::cas[]
	getRes, err := collection.Get("player432", &gocb.GetOptions{})
	if err != nil {
		panic(err)
	}

	mops = []gocb.MutateInSpec{
		gocb.DecrementSpec("gold", 150, &gocb.CounterSpecOptions{}),
	}
	decrementCasResult, err := collection.MutateIn("player432", mops, &gocb.MutateInOptions{
		Cas: getRes.Cas(),
	})
	if err != nil {
		panic(err)
	}
	// #end::cas[]
	fmt.Println(decrementCasResult)
}
