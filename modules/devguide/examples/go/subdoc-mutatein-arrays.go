package main

import (
	"errors"
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

	// Array Append
	// #tag::mutateInArrayAppend[]
	mops := []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("purchases.complete", 777, nil),
	}
	arrayAppendResult, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		panic(err)
	}
	// purchases.complete is now [339, 976, 442, 666, 777]
	// #end::mutateInArrayAppend[]
	fmt.Println(arrayAppendResult.Cas())

	// Array Prepend
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
	fmt.Println(arrayPrependResult.Cas())

	// Array Doc
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
	fmt.Println(arrayAppendRootResult.Cas)

	// Array Multiples
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
	fmt.Println(multiArrayAppendResult.Cas())

	// Array Multiples as one element
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
	fmt.Println(singleArrayAppendResult.Cas())

	// Array multiple specs
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
	fmt.Println(individualArrayAppendResult.Cas())

	// Array Create document path
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
	fmt.Println(createPathResult.Cas())

	// Array Add Unique
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
	fmt.Println(arrayAddUniqueResult.Cas())
	fmt.Println(arrayAddUniqueSecondResult.Cas())

	// Array Add Insert
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
	fmt.Println(arrayInsertResult.Cas())

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
