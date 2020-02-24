package main

import (
	"fmt"

	gocb "github.com/couchbase/gocb/v2"
)

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

	bucket := cluster.Bucket("bucket-name")

	collection := bucket.DefaultCollection()

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
	fmt.Println(firstConcurrentResult)
	fmt.Println(secondConcurrentResult)

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
