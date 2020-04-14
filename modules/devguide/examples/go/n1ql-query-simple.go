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

	// #tag::results[]
	query := "SELECT x.* FROM `travel-sample` x WHERE x.`type`=$1 LIMIT 10;"
	rows, err := cluster.Query(query, &gocb.QueryOptions{PositionalParameters: []interface{}{"hotel"}})
	// check query was successful
	if err != nil {
		panic(err)
	}

	type hotel struct {
		Name string `json:"name"`
	}

	var hotels []hotel
	// iterate over rows
	for rows.Next() {
		var h hotel // this could also just be an interface{} type
		err := rows.Row(&h)
		if err != nil {
			panic(err)
		}
		hotels = append(hotels, h)
	}

	// always check for errors after iterating
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	// #end::results[]

	// #tag::metrics[]
	metadata, err := rows.MetaData()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Execution Time: %d\n", metadata.Metrics.ExecutionTime)
	// #end::metrics[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
