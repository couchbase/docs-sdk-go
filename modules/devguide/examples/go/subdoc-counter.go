package main

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.Connect("10.112.194.101", opts)
	if err != nil {
		panic(err)
	}

	collection := cluster.Bucket("travel-sample").DefaultCollection()

	// Increment
	// #tag::mutateInIncrement[]
	mops := []gocb.MutateInSpec{
		gocb.IncrementSpec("logins", 1, &gocb.CounterSpecOptions{}),
	}
	incrementResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}

	var logins int
	err = incrementResult.ContentAt(0, &logins)
	if err != nil {
		panic(err)
	}
	fmt.Println(logins) // 1
	// #end::mutateInIncrement[]

	// Decrement
	// #tag::mutateInDecrement[]
	_, err = collection.Upsert("player432", map[string]int{"gold": 1000}, &gocb.UpsertOptions{})
	if err != nil {
		panic(err)
	}

	mops = []gocb.MutateInSpec{
		gocb.DecrementSpec("gold", 150, &gocb.CounterSpecOptions{}),
	}
	decrementResult, err := collection.MutateIn("player432", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}

	var gold int
	err = decrementResult.ContentAt(0, &gold)
	if err != nil {
		panic(err)
	}
	fmt.Printf("player 432 now has %d gold remaining\n", gold)
	// player 432 now has 850 gold remaining
	// #end::mutateInDecrement[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
