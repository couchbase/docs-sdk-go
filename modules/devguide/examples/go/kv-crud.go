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
	cluster, err := gocb.Connect("your-ip", opts)
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

	// #tag::insert[]
	// Insert Document
	type myDoc struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	document := myDoc{Foo: "bar", Bar: "foo"}
	result, err := collection.Insert("document-key", &document, nil)
	if err != nil {
		panic(err)
	}
	// #end::insert[]
	fmt.Println(result)

	// #tag::insertoptions[]
	// Insert Document with options
	resultwithOptions, err := collection.Insert("document-key-options", &document, &gocb.InsertOptions{
		Timeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	// #end::insertoptions[]
	fmt.Println(resultwithOptions)

	// #tag::replacecas[]
	// Replace Document with Cas
	replaceResultWithCas, err := collection.Replace("document-key", &document, &gocb.ReplaceOptions{
		Cas: 12345,
	})
	if err != nil {
		// We expect this to error
		fmt.Println(err)
	}
	// #end::replacecas[]
	fmt.Println(replaceResultWithCas)

	// #tag::update[]
	// Get and Replace Document with Cas
	updateGetResult, err := collection.Get("document-key", nil)
	if err != nil {
		panic(err)
	}

	var doc myDoc
	err = updateGetResult.Content(&doc)
	if err != nil {
		panic(err)
	}

	doc.Bar = "moo"

	updateResult, err := collection.Replace("document-key", doc, &gocb.ReplaceOptions{
		Cas: updateGetResult.Cas(),
	})
	// #end::update[]
	fmt.Println(updateResult)

	// #tag::get[]
	// Get
	getResult, err := collection.Get("document-key", nil)
	if err != nil {
		panic(err)
	}

	var getDoc myDoc
	err = getResult.Content(&getDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(getDoc)
	// #end::get[]

	// #tag::gettimeout[]
	// Get with timeout
	getTimeoutResult, err := collection.Get("document-key", &gocb.GetOptions{
		Timeout: 10 * time.Millisecond,
	})
	if err != nil {
		panic(err)
	}

	var getTimeoutDoc myDoc
	err = getTimeoutResult.Content(&getTimeoutDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(getTimeoutDoc)
	// #end::gettimeout[]

	// #tag::remove[]
	// Remove with Durability
	removeResult, err := collection.Remove("document-key", &gocb.RemoveOptions{
		Timeout:         100 * time.Millisecond,
		DurabilityLevel: gocb.DurabilityLevelMajority,
	})
	if err != nil {
		panic(err)
	}
	// #end::remove[]
	fmt.Println(removeResult)

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
