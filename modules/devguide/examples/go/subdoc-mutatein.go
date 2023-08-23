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
	cluster, err := gocb.Connect("your-ip", opts)
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

	// Upsert
	// #tag::mutateInUpsert[]
	mops := []gocb.MutateInSpec{
		gocb.UpsertSpec("fax", "311-555-0151", &gocb.UpsertSpecOptions{}),
	}
	upsertResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{
		Timeout: 10050 * time.Millisecond,
	})
	if err != nil {
		panic(err)
	}
	// #end::mutateInUpsert[]
	fmt.Println(upsertResult.Cas())

	// Insert
	// #tag::mutateInInsert[]
	mops = []gocb.MutateInSpec{
		gocb.InsertSpec("purchases.pending", []interface{}{32, true, "None"}, &gocb.InsertSpecOptions{}),
	}
	insertResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// #end::mutateInInsert[]
	fmt.Println(insertResult.Cas())

	// Multiple specs
	// #tag::mutateInMulti[]
	mops = []gocb.MutateInSpec{
		gocb.RemoveSpec("addresses.billing", nil),
		gocb.ReplaceSpec("email", "dougr96@hotmail.com", nil),
	}
	multiMutateResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// #end::mutateInMulti[]
	fmt.Println(multiMutateResult.Cas())

	// Create path
	// #tag::mutateInCreatePath[]
	mops = []gocb.MutateInSpec{
		gocb.UpsertSpec("level_0.level_1.foo.bar.phone", map[string]interface{}{
			"num": "311-555-0101",
			"ext": 16,
		}, &gocb.UpsertSpecOptions{
			CreatePath: true,
		}),
	}
	createPathUpsertResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// #end::mutateInCreatePath[]
	fmt.Println(createPathUpsertResult.Cas())

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
