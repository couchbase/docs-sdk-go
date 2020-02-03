package main

import (
	"errors"
	"fmt"
	"time"

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

	// #tag::mutateInUpsert[]
	mops := []gocb.MutateInSpec{
		gocb.UpsertSpec("fax", "311-555-0151", &gocb.UpsertSpecOptions{}),
	}
	upsertResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{
		Timeout: 50 * time.Millisecond,
	})
	if err != nil {
		panic(err)
	}
	// #end::mutateInUpsert[]
	fmt.Println(upsertResult)

	// #tag::mutateInInsert[]
	mops = []gocb.MutateInSpec{
		gocb.InsertSpec("purchases.complete", []interface{}{32, true, "None"}, &gocb.InsertSpecOptions{}),
	}
	insertResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// #end::mutateInInsert[]
	fmt.Println(insertResult)

	// #tag::mutateInMulti[]
	mops = []gocb.MutateInSpec{
		gocb.RemoveSpec("addresses.billing[2]", nil),
		gocb.ReplaceSpec("email", "dougr96@hotmail.com", nil),
	}
	multiMutateResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// #end::mutateInMulti[]
	fmt.Println(multiMutateResult)

	// #tag::mutateInArrayAppend[]
	mops = []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("purchases.complete", 777, nil),
	}
	arrayAppendResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// purchases.complete is now [339, 976, 442, 666, 777]
	// #end::mutateInArrayAppend[]
	fmt.Println(arrayAppendResult)

	// #tag::mutateInArrayPrepend[]
	mops = []gocb.MutateInSpec{
		gocb.ArrayPrependSpec("purchases.abandoned", 17, nil),
	}
	arrayPrependResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// purchases.abandoned is now [18, 157, 49, 999]
	// #end::mutateInArrayPrepend[]
	fmt.Println(arrayPrependResult)

	// #tag::mutateInArrayDoc[]
	upsertDocResult, err := collection.Upsert("my_array", []int{}, nil)
	if err != nil {
		panic(err)
	}

	mops = []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("", "some element", &gocb.ArrayAppendSpecOptions{}),
	}
	arrayAppendRootResult, err := collection.MutateIn("my_array", mops, &gocb.MutateInOptions{
		Cas: upsertDocResult.Cas(),
	})
	if err != nil {
		panic(err)
	}
	// the document my_array is now ["some element"]
	// #end::mutateInArrayDoc[]
	fmt.Println(arrayAppendRootResult)

	// #tag::mutateInArrayDocMulti[]
	mops = []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("", []string{"elem1", "elem2", "elem3"}, &gocb.ArrayAppendSpecOptions{
			HasMultiple: true, // this signifies that the value is multiple array elements rather than one
		}),
	}
	multiArrayAppendResult, err := collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// the document my_array is now ["some_element", "elem1", "elem2", "elem3"]
	// #end::mutateInArrayDocMulti[]
	fmt.Println(multiArrayAppendResult)

	// #tag::mutateInArrayDocMultiSingle[]
	mops = []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("", []string{"elem1", "elem2", "elem3"}, &gocb.ArrayAppendSpecOptions{}),
	}
	singleArrayAppendResult, err := collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// the document my_array is now ["some_element", "elem1", "elem2", "elem3", ["elem1", "elem2", "elem3"]]
	// #end::mutateInArrayDocMultiSingle[]
	fmt.Println(singleArrayAppendResult)

	// #tag::mutateInArrayAppendMulti[]
	mops = []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("", "elem1", &gocb.ArrayAppendSpecOptions{}),
		gocb.ArrayAppendSpec("", "elem2", &gocb.ArrayAppendSpecOptions{}),
		gocb.ArrayAppendSpec("", "elem3", &gocb.ArrayAppendSpecOptions{}),
	}
	individualArrayAppendResult, err := collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// #end::mutateInArrayAppendMulti[]
	fmt.Println(individualArrayAppendResult)

	// #tag::mutateInArrayAppendCreatePath[]
	mops = []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("some.array", []string{"Hello", "World"}, &gocb.ArrayAppendSpecOptions{
			HasMultiple: true,
			CreatePath:  true,
		}),
	}
	createPathResult, err := collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// #end::mutateInArrayAppendCreatePath[]
	fmt.Println(createPathResult)

	// #tag::mutateInArrayAddUnique[]
	mops = []gocb.MutateInSpec{
		gocb.ArrayAddUniqueSpec("purchases.complete", 95, &gocb.ArrayAddUniqueSpecOptions{}),
	}
	arrayAddUniqueResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}

	mops = []gocb.MutateInSpec{
		gocb.ArrayAddUniqueSpec("purchases.complete", 95, &gocb.ArrayAddUniqueSpecOptions{}),
	}
	arrayAddUniqueSecondResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	fmt.Println(errors.Is(err, gocb.ErrPathExists)) // true
	// #end::mutateInArrayAddUnique[]
	fmt.Println(arrayAddUniqueResult)
	fmt.Println(arrayAddUniqueSecondResult)

	// #tag::mutateInArrayInsert[]
	mops = []gocb.MutateInSpec{
		gocb.ArrayInsertSpec("some.array[1]", "Cruel", &gocb.ArrayInsertSpecOptions{}),
	}
	arrayInsertResult, err := collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// The value at some.array is now [Hello, Cruel, World]
	// #end::mutateInArrayInsert[]
	fmt.Println(arrayInsertResult)

	// #tag::mutateInIncrement[]
	mops = []gocb.MutateInSpec{
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
	fmt.Println(createPathUpsertResult)

	// #tag::concurrent[]
	mops = []gocb.MutateInSpec{
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

	// #tag::traddurability[]
	mops = []gocb.MutateInSpec{
		gocb.InsertSpec("name", "mike", nil),
	}
	observeResult, err := collection.MutateIn("key", mops, &gocb.MutateInOptions{
		PersistTo:   1,
		ReplicateTo: 1,
	})
	if err != nil {
		panic(err)
	}
	// #end::traddurability[]
	fmt.Println(observeResult)

	// #tag::newdurability[]
	mops = []gocb.MutateInSpec{
		gocb.InsertSpec("name", "mike", nil),
	}
	durableResult, err := collection.MutateIn("key", mops, &gocb.MutateInOptions{
		DurabilityLevel: gocb.DurabilityLevelMajority,
	})
	if err != nil {
		panic(err)
	}
	// #end::newdurability[]
	fmt.Println(durableResult)
}
