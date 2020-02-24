package main

import (
	"fmt"
	"time"

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

	// Get
	// #tag::lookupInGet[]
	ops := []gocb.LookupInSpec{
		gocb.GetSpec("addresses.delivery.country", &gocb.GetSpecOptions{}),
	}
	getResult, err := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{})
	if err != nil {
		panic(err)
	}

	var country string
	err = getResult.ContentAt(0, &country)
	if err != nil {
		panic(err)
	}
	fmt.Println(country) // United Kingdom
	// #end::lookupInGet[]

	// Exists
	// #tag::lookupInExists[]
	ops = []gocb.LookupInSpec{
		gocb.ExistsSpec("purchases.pending[-1]", &gocb.ExistsSpecOptions{}),
	}
	existsResult, err := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{})
	if err != nil {
		panic(err)
	}

	var exists bool
	err = existsResult.ContentAt(0, &exists)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Path exists? %t\n", exists) // Path exists? false
	// #end::lookupInExists[]

	// Multiple specs
	// #tag::lookupInMulti[]
	ops = []gocb.LookupInSpec{
		gocb.GetSpec("addresses.delivery.country", nil),
		gocb.ExistsSpec("purchases.pending[-1]", nil),
	}
	multiLookupResult, err := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{
		Timeout: 50 * time.Millisecond,
	})
	if err != nil {
		panic(err)
	}

	var multiCountry string
	err = multiLookupResult.ContentAt(0, &multiCountry)
	if err != nil {
		panic(err)
	}
	var multiExists bool
	err = multiLookupResult.ContentAt(1, &multiExists)
	if err != nil {
		panic(err)
	}

	fmt.Println(multiCountry)                    // United Kingdom
	fmt.Printf("Path exists? %t\n", multiExists) // Path exists? false
	// #end::lookupInMulti[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
