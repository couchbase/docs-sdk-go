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

	// #tag::one[]
	query := "SELECT x.* FROM `travel-sample` x WHERE x.`type`=$1 LIMIT 1;"
	rows, err := cluster.Query(query, &gocb.QueryOptions{PositionalParameters: []interface{}{"hotel"}})

	// check query was successful
	if err != nil {
		panic(err)
	}

	var hotel interface{} // this could also be a specific type like Hotel
	err = rows.One(&hotel)
	if err != nil {
		panic(err)
	}
	fmt.Println(hotel)
	// #end::one[]

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
