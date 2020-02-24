package main

import (
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

	collection := cluster.Bucket("bucket-name").DefaultCollection()

	// #tag::consistentwith[]
	// create / update document (mutation)
	result, err := collection.Upsert("id", struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}{Name: "somehotel", Type: "hotel"}, nil)
	if err != nil {
		panic(err)
	}

	// create mutation state from mutation results
	state := gocb.NewMutationState(*result.MutationToken())

	// use mutation state with query option
	rows, err := cluster.Query("SELECT x.* FROM `travel-sample` x WHERE x.`type`=\"hotel\" LIMIT 10", &gocb.QueryOptions{
		ConsistentWith: state,
	})
	// #end::consistentwith[]
	if err != nil {
		panic(err)
	}

	// iterate over rows
	for rows.Next() {
		var hotel interface{} // this could also be a specific type like Hotel
		err := rows.Row(&hotel)
		if err != nil {
			panic(err)
		}
		fmt.Println(hotel)
	}

	// always check for errors after iterating
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	metadata, err := rows.MetaData()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Execution Time: %d\n", metadata.Metrics.ExecutionTime)

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
