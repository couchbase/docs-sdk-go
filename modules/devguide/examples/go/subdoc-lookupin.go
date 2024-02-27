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
	collection := bucket.Scope("inventory").Collection("airline")

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

	// Get
	{
		// #tag::lookupInGet[]
		ops := []gocb.LookupInSpec{
			gocb.GetSpec("addresses.delivery.country", &gocb.GetSpecOptions{}),
		}
		res, err := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{})
		if err != nil {
			panic(err)
		}

		var country string
		err = res.ContentAt(0, &country)
		if err != nil {
			panic(err)
		}
		fmt.Println(country) // United Kingdom
		// #end::lookupInGet[]
	}

	// Exists
	{
		// #tag::lookupInExists[]
		ops := []gocb.LookupInSpec{
			gocb.ExistsSpec("purchases.pending[-1]", &gocb.ExistsSpecOptions{}),
		}
		res, err := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{})
		if err != nil {
			panic(err)
		}

		exists := res.Exists(0)

		fmt.Printf("Path exists? %t\n", exists) // Path exists? false
		// #end::lookupInExists[]
	}

	// Multiple specs
	{
		// #tag::lookupInMulti[]
		ops := []gocb.LookupInSpec{
			gocb.GetSpec("addresses.delivery.country", nil),
			gocb.ExistsSpec("purchases.pending[-1]", nil),
		}
		res, err := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{
			Timeout: 50 * time.Millisecond,
		})
		if err != nil {
			panic(err)
		}

		var country string
		err = res.ContentAt(0, &country)
		if err != nil {
			panic(err)
		}
		exists := res.Exists(1)

		fmt.Println(country)                    // United Kingdom
		fmt.Printf("Path exists? %t\n", exists) // Path exists? false
		// #end::lookupInMulti[]
	}

	// From any replica
	{
		// #tag::lookupInAnyReplica[]
		ops := []gocb.LookupInSpec{
			gocb.GetSpec("addresses.delivery.country", &gocb.GetSpecOptions{}),
		}
		res, err := collection.LookupInAnyReplica("customer123", ops, nil)
		if err != nil {
			panic(err)
		}

		var country string
		err = res.ContentAt(0, &country)
		if err != nil {
			panic(err)
		}

		fmt.Println(country) // United Kingdom
		fmt.Printf("Is replica? %t\n", res.IsReplica())
		// #end::lookupInAnyReplica[]
	}

	// From all replicas
	{
		// #tag::lookupInAllReplicas[]
		ops := []gocb.LookupInSpec{
			gocb.GetSpec("addresses.delivery.country", &gocb.GetSpecOptions{}),
		}
		stream, err := collection.LookupInAllReplicas("customer123", ops, nil)
		if err != nil {
			panic(err)
		}

		for {
			replicaRes := stream.Next()
			if replicaRes == nil {
				break
			}

			var country string
			err = replicaRes.ContentAt(0, &country)
			if err != nil {
				panic(err)
			}

			fmt.Println(country) // United Kingdom
			fmt.Printf("Is replica? %t\n", replicaRes.IsReplica())
		}
		// #end::lookupInAllReplicas[]
	}

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
