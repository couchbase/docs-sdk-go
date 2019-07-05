package main

// #tag::connect[]
import (
	"fmt"

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
		// handle err
	}
	// #end::connect[]

	// #tag::bucket[]
	// get a bucket reference
	bucket := cluster.Bucket("bucket-name", &gocb.BucketOptions{})
	// #end::bucket[]

	// #tag::collection[]
	// get a collection reference
	collection := bucket.DefaultCollection(&gocb.CollectionOptions{})
	// for a named collection and scope
	// collection := bucket.Scope("my-scope").Collection("my-collection", &gocb.CollectionOptions{})
	// #end::collection[]

	// #tag::lookupInGet[]
	spec := gocb.LookupInSpec{}
	ops := []gocb.LookupInOp{
		spec.Get("addresses.delivery.country", &gocb.LookupInSpecGetOptions{}),
	}
	getResult, err := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{})
	if err != nil {
		panic(err)
	}

	var country string
	err = getResult.ContentAt(0, &country)
	if err != nil {
		// handle err
	}
	fmt.Println(country) // United Kingdom
	// #end::lookupInGet[]

	// #tag::lookupInExists[]
	ops = []gocb.LookupInOp{
		spec.Exists("purchases.pending[-1]", &gocb.LookupInSpecExistsOptions{}),
	}
	existsResult, _ := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{})
	var exists bool
	existsResult.ContentAt(0, &exists)

	fmt.Printf("Path exists? %t\n", exists) // Path exists? false
	// #end::lookupInExists[]

	// #tag::lookupInMulti[]
	ops = []gocb.LookupInOp{
		spec.Get("addresses.delivery.country", &gocb.LookupInSpecGetOptions{}),
		spec.Exists("purchases.pending[-1]", &gocb.LookupInSpecExistsOptions{}),
	}
	multiLookupResult, _ := collection.LookupIn("customer123", ops, &gocb.LookupInOptions{})

	var multiCountry string
	err = multiLookupResult.ContentAt(0, &multiCountry)
	if err != nil {
		// handle err
	}
	var multiExists bool
	multiLookupResult.ContentAt(1, &multiExists)

	fmt.Println(multiCountry)                    // United Kingdom
	fmt.Printf("Path exists? %t\n", multiExists) // Path exists? false
	// #end::lookupInMulti[]

	// #tag::mutateInUpsert[]
	mSpec := gocb.MutateInSpec{}
	mops := []gocb.MutateInOp{
		mSpec.Upsert("fax", "311-555-0151", &gocb.MutateInSpecUpsertOptions{}),
	}
	_, err = collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	if err != nil {
		// handle err
	}
	// #end::mutateInUpsert[]

	// #tag::mutateInInsert[]
	mops = []gocb.MutateInOp{
		mSpec.Insert("purchases.complete", []interface{}{32, true, "None"}, &gocb.MutateInSpecInsertOptions{}),
	}
	collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	// #end::mutateInInsert[]

	// #tag::mutateInMulti[]
	mops = []gocb.MutateInOp{
		mSpec.Remove("addresses.billing", &gocb.MutateInSpecRemoveOptions{}),
		mSpec.Replace("email", "dougr96@hotmail.com", &gocb.MutateInSpecReplaceOptions{}),
	}
	collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	// #end::mutateInMulti[]

	// #tag::mutateInArrayAppend[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayAppend("purchases.complete", 777, &gocb.MutateInSpecArrayAppendOptions{}),
	}
	collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	// purchases.complete is now [339, 976, 442, 666, 777]
	// #end::mutateInArrayAppend[]

	// #tag::mutateInArrayPrepend[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayPrepend("purchases.abandoned", 17, &gocb.MutateInSpecArrayPrependOptions{}),
	}
	collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	// purchases.abandoned is now [18, 157, 49, 999]
	// #end::mutateInArrayPrepend[]

	// #tag::mutateInArrayDoc[]
	collection.Upsert("my_array", []int{}, &gocb.UpsertOptions{})
	mops = []gocb.MutateInOp{
		mSpec.ArrayAppend("", "some element", &gocb.MutateInSpecArrayAppendOptions{}),
	}
	collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	// the document my_array is now ["some element"]
	// #end::mutateInArrayDoc[]

	// #tag::mutateInArrayDocMulti[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayAppend("", []string{"elem1", "elem2", "elem3"}, &gocb.MutateInSpecArrayAppendOptions{
			HasMultiple: true, // this signifies that the value is multiple array elements rather than one
		}),
	}
	collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	// the document my_array is now ["some_element", "elem1", "elem2", "elem3"]
	// #end::mutateInArrayDocMulti[]

	// #tag::mutateInArrayDocMultiSingle[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayAppend("", []string{"elem1", "elem2", "elem3"}, &gocb.MutateInSpecArrayAppendOptions{}),
	}
	collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	// the document my_array is now ["some_element", "elem1", "elem2", "elem3", ["elem1", "elem2", "elem3"]]
	// #end::mutateInArrayDocMultiSingle[]

	// #tag::mutateInArrayAppendMulti[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayAppend("", "elem1", &gocb.MutateInSpecArrayAppendOptions{}),
		mSpec.ArrayAppend("", "elem2", &gocb.MutateInSpecArrayAppendOptions{}),
		mSpec.ArrayAppend("", "elem3", &gocb.MutateInSpecArrayAppendOptions{}),
	}
	collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	// #end::mutateInArrayAppendMulti[]

	// #tag::mutateInArrayAppendCreatePath[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayAppend("some.array", []string{"Hello", "World"}, &gocb.MutateInSpecArrayAppendOptions{
			HasMultiple: true,
			CreatePath:  true,
		}),
	}
	collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	// #end::mutateInArrayAppendCreatePath[]

	// #tag::mutateInArrayAddUnique[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayAddUnique("purchases.complete", 95, &gocb.MutateInSpecArrayAddUniqueOptions{}),
	}
	_, err = collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	// err is nil

	mops = []gocb.MutateInOp{
		mSpec.ArrayAddUnique("purchases.complete", 95, &gocb.MutateInSpecArrayAddUniqueOptions{}),
	}
	_, err = collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	fmt.Println(gocb.IsPathExistsError(err)) // true
	// #end::mutateInArrayAddUnique[]

	// #tag::mutateInArrayInsert[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayInsert("some.array[1]", "Cruel", &gocb.MutateInSpecArrayInsertOptions{}),
	}
	collection.MutateIn("my_array", mops, &gocb.MutateInOptions{})
	// #end::mutateInArrayInsert[]

	// #tag::mutateInIncrement[]
	mops = []gocb.MutateInOp{
		mSpec.Increment("logins", 1, &gocb.MutateInSpecCounterOptions{}),
	}
	incrementResult, _ := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})

	var logins int
	incrementResult.ContentAt(0, &logins)
	fmt.Println(logins) // 1
	// #end::mutateInIncrement[]

	// #tag::mutateInDecrement[]
	collection.Upsert("player432", map[string]int{"gold": 1000}, &gocb.UpsertOptions{})
	mops = []gocb.MutateInOp{
		mSpec.Decrement("gold", 150, &gocb.MutateInSpecCounterOptions{}),
	}
	decrementResult, _ := collection.MutateIn("player432", mops, &gocb.MutateInOptions{})

	var gold int
	decrementResult.ContentAt(0, &gold)
	fmt.Printf("player 432 now has %d gold remaining\n", gold)
	// player 432 now has 850 gold remaining
	// #end::mutateInDecrement[]

	// #tag::mutateInCreatePath[]
	mops = []gocb.MutateInOp{
		mSpec.Upsert("level_0.level_1.foo.bar.phone", map[string]interface{}{
			"num": "311-555-0101",
			"ext": 16,
		}, &gocb.MutateInSpecUpsertOptions{
			CreatePath: true,
		}),
	}
	collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	// #end::mutateInCreatePath[]

	// #tag::concurrent[]
	mops = []gocb.MutateInOp{
		mSpec.ArrayAppend("purchases.complete", 99, &gocb.MutateInSpecArrayAppendOptions{}),
	}
	collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})

	mops = []gocb.MutateInOp{
		mSpec.ArrayAppend("purchases.abandoned", 101, &gocb.MutateInSpecArrayAppendOptions{}),
	}
	collection.MutateIn("customer123", mops, &gocb.MutateInOptions{})
	// #end::concurrent[]

	// #tag::cas[]
	getRes, err := collection.Get("player432", &gocb.GetOptions{})
	if err != nil {
		// handle error
	}

	mops = []gocb.MutateInOp{
		mSpec.Decrement("gold", 150, &gocb.MutateInSpecCounterOptions{}),
	}
	collection.MutateIn("player432", mops, &gocb.MutateInOptions{
		Cas: getRes.Cas(),
	})
	// #end::cas[]

	// #tag::traddurability[]
	mops = []gocb.MutateInOp{
		mSpec.Insert("name", "mike", nil),
	}
	collection.MutateIn("key", mops, &gocb.MutateInOptions{
		PersistTo:   1,
		ReplicateTo: 1,
	})
	// #end::traddurability[]

	// #tag::newdurability[]
	mops = []gocb.MutateInOp{
		mSpec.Insert("name", "mike", nil),
	}
	collection.MutateIn("key", mops, &gocb.MutateInOptions{
		DurabilityLevel: gocb.DurabilityLevelMajority,
	})
	// #end::newdurability[]
}
