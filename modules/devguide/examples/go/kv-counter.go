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

	// Increment & Decrement are considered part of the 'binary' API and as such may still be subject to change.

	// Create a document and assign it to 10 - counter works atomically by
	// first creating a document if it doesn't exist.   If it exists, the
	// same method will increment/decrement per the "delta" parameter
	// #tag::increment[]
	binaryC := collection.Binary()
	key := "goDevguideExampleCounter"
	curKeyValue, err := binaryC.Increment(key, &gocb.IncrementOptions{
		Initial: 10,
		Delta:   2,
	})
	if err != nil {
		panic(err)
	}
	// #end::increment[]

	// Should Print 10
	fmt.Println("Initialized Counter:", curKeyValue)

	// Issue same operation, increment value by 2, to 12
	curKeyValue, err = binaryC.Increment(key, &gocb.IncrementOptions{
		Initial: 10,
		Delta:   2,
	})
	if err != nil {
		panic(err)
	}

	// Should Print 12
	fmt.Println("Incremented Counter:", curKeyValue)

	// #tag::decrement[]
	// Issue same operation, increment value by 2, to 12
	curKeyValue, err = binaryC.Decrement(key, &gocb.DecrementOptions{
		Initial: 10,
		Delta:   4,
	})
	if err != nil {
		panic(err)
	}
	// #end::decrement[]

	// Should Print 8
	fmt.Println("Decremented Counter:", curKeyValue)

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
