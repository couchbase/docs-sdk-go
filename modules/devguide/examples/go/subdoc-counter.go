package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	bucket := cluster.Bucket("travel-sample")
	collection := bucket.DefaultCollection()

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	var customer123 interface{}
	b, err := ioutil.ReadFile("customer123.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &customer123)
	if err != nil {
		panic(err)
	}

	_, err = collection.Upsert("customer123", customer123, nil)
	if err != nil {
		panic(err)
	}

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
