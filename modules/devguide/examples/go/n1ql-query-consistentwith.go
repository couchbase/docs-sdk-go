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

	// For Server versions 6.5 or later you do not need to open a bucket here
	b := cluster.Bucket("travel-sample")

	// We wait until the bucket is definitely connected and setup.
	// For Server versions 6.5 or later if we hadn't opened a bucket then we could use cluster.WaitUntilReady here.
	err = b.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	collection := b.Scope("inventory").Collection("hotel")

	// NOTE: This currently fails with Couchbase Internal Server error.
	// Server issue tracked here: https://issues.couchbase.com/browse/MB-46876
	// Add back in once Couchbase Server 7.0.1 is available, which will fix this issue.
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
	rows, err := cluster.Query(
		"SELECT x.* FROM `travel-sample`.inventory.hotel x WHERE x.`city`= $1 LIMIT 10",
		&gocb.QueryOptions{
			ConsistentWith:       state,
			PositionalParameters: []interface{}{"San Francisco"},
			Adhoc:                true,
		},
	)
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

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
