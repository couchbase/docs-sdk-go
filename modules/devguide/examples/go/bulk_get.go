package main

import (
	"fmt"
	"log"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func Example_concurrentGet() {
	// #tag::connect[]
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

	// We wait until the bucket is connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}
	// #end::connect[]

	// #tag::prepareOps[]
	// Replace with your actual document IDs.
	docIDs := []string{
		"airline_10",
		"airline_10123",
		"airline_10226",
		"airline_10642",
		"airline_10748",
		"airline_10765",
		"airline_109",
		"airline_792",
		"airline_8745",
		"airline_8809",
		"airline_9833",
	}

  	// #tag::bulk-get[]
	var getOps []gocb.BulkOp
	for _, id := range docIDs {
		getOps = append(getOps, &gocb.GetOp{
			ID: id,
		})
	}
	// #end::bulk-get[]
  // #end::prepareOps[]

	// #tag::send[]
	err = collection.Do(getOps, nil)
	if err != nil {
		log.Println(err)
	}

	// Be sure to check each individual operation for errors too.
	var fetchedCount int
	for _, op := range getOps {
		getOp, ok := op.(*gocb.GetOp)
		if !ok {
			log.Println("failed to type assert BulkOp to *gocb.GetOp")
			continue
		}

		if getOp.Err != nil {
			log.Printf("failed to fetch document %s: %v\n", getOp.ID, getOp.Err)
			continue
		}

		var docContent interface{}
		if err := getOp.Result.Content(&docContent); err != nil {
			log.Printf("failed to decode document %s content: %v\n", getOp.ID, err)
			continue
		}

		fetchedCount++
		fmt.Printf("Fetched document #%d: ID: %s, Content: %v\n", fetchedCount, getOp.ID, docContent)
	}

	log.Printf("Completed fetching %d/%d documents\n", fetchedCount, len(docIDs))
	// #end::send[]

	cluster.Close(nil)
}
