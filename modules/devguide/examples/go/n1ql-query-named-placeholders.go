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

	// #tag::named-params[]
	query := "SELECT x.* FROM `travel-sample` x WHERE x.`type`=$type LIMIT 10;"
	params := make(map[string]interface{}, 1)
	params["type"] = "hotel"
	rows, err := cluster.Query(query, &gocb.QueryOptions{NamedParameters: params})
	// #end::named-params[]

	// check query was successful
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
