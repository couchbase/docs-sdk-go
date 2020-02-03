package main

// #tag::connect[]
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
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("bucket-name")

	collection := bucket.DefaultCollection()

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
	resultwithOptions, err := collection.Insert("document-key", &document, &gocb.InsertOptions{
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
		panic(err)
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

	// #tag::expiry[]
	// Upsert with Expiry
	expiryResult, err := collection.Upsert("document-key", &document, &gocb.UpsertOptions{
		Timeout: 25 * time.Millisecond,
		Expiry:  60, // Seconds
	})
	// #end::expiry[]
	fmt.Println(expiryResult)

	// #tag::durability[]
	// Upsert with Durability
	durableResult, err := collection.Upsert("document-key", &document, &gocb.UpsertOptions{
		DurabilityLevel: gocb.DurabilityLevelMajority,
	})
	// #end::durability[]
	fmt.Println(durableResult)

	// #tag::observebased[]
	// Upsert with Observe based durability
	observeResult, err := collection.Upsert("document-key", &document, &gocb.UpsertOptions{
		PersistTo:   1, // Has been written to disk on 1 other node than active
		ReplicateTo: 1, // Has been written to memory on 1 other node than active
	})
	// #end::observebased[]
	fmt.Println(observeResult)

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

	// #tag::touch[]
	// Touch
	touchResult, err := collection.Touch("document-key", 60, &gocb.TouchOptions{
		Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		panic(err)
	}
	// #end::touch[]
	fmt.Println(touchResult)

	// #tag::getandtouch[]
	// GetAndTouch
	getAndTouchResult, err := collection.GetAndTouch("document-key", 60, &gocb.GetAndTouchOptions{
		Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		panic(err)
	}

	var getAndTouchDoc myDoc
	err = getAndTouchResult.Content(&getAndTouchDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(getAndTouchDoc)
	// #end::getandtouch[]

	// #tag::increment[]
	// Increment
	incrementResult, err := collection.Binary().Increment("document-key", &gocb.IncrementOptions{
		Initial: 1000,
		Delta:   1,
		Timeout: 50 * time.Millisecond,
		Expiry:  3600, // Seconds
	})
	if err != nil {
		panic(err)
	}
	// #end::increment[]
	fmt.Println(incrementResult)

	// #tag::decrement[]
	// Increment
	decrementResult, err := collection.Binary().Decrement("document-key", &gocb.DecrementOptions{
		Initial: 1000,
		Delta:   1,
		Timeout: 50 * time.Millisecond,
		Expiry:  3600, // Seconds
	})
	if err != nil {
		panic(err)
	}
	// #end::decrement[]
	fmt.Println(decrementResult)
}
