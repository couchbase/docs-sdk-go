package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// Connect to Cluster
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

	// Open Bucket and collection
	bucket := cluster.Bucket("default")
	collection := bucket.DefaultCollection()

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	// Create a document and assign a "Persist To" value of 1 node.
	// Should Always Succeed, even on single node cluster.
	key := "goDevguideExamplePersistTo"
	val := "Durabilty PersistTo Test Value"
	_, err = collection.Upsert(key, &val, &gocb.UpsertOptions{
		PersistTo: 1,
	})
	if err != nil {
		panic(err)
	}

	// Retrieve Value Persist To
	getRes, err := collection.Get(key, nil)
	if err != nil {
		panic(err)
	}

	var retValue interface{}
	err = getRes.Content(&retValue)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document Retrieved:", retValue)

	// Create a document and assign a "Replicate To" value of 1 node.
	// Should Fail on a single node cluster, succeed on a multi node
	// cluster of 3 or more nodes with at least one replica enabled.
	key = "goDevguideExampleReplicateTo"
	val = "Durabilty ReplicateTo Test Value"
	_, err = collection.Upsert(key, &val, &gocb.UpsertOptions{
		ReplicateTo: 1,
	})
	if err != nil {
		panic(err)
	}

	// Retrieve Value Replicate To
	// Should succeed even if durability fails, as the document was
	// still written.
	getRes, err = collection.Get(key, nil)
	if err != nil {
		panic(err)
	}

	err = getRes.Content(&retValue)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document Retrieved:", retValue)

	// Create a document and assign a "Replicate To" and a "Persist TO"
	// value of 1 node. Should Fail on a single node cluster, succeed on
	// a multi node cluster of 3 or more nodes with at least one replica
	// enabled.
	// #tag::observebased[]
	key = "replicateToAndPersistTo"
	val = "Durabilty ReplicateTo and PersistTo Test Value"
	_, err = collection.Upsert(key, &val, &gocb.UpsertOptions{
		PersistTo:   1,
		ReplicateTo: 1,
	})
	if err != nil {
		panic(err)
	}
	// #end::observebased[]

	// Retrieve Value Replicate To and Persist To
	// Should succeed even if durability fails, as the document was
	// still written.
	getRes, err = collection.Get(key, nil)
	if err != nil {
		panic(err)
	}

	err = getRes.Content(&retValue)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document Retrieved:", retValue)

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
