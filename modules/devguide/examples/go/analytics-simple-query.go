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

	// #tag::simple[]
	results, err := cluster.AnalyticsQuery("select airportname, country from airports where country = 'France';", nil)
	if err != nil {
		panic(err)
	}
	// #end::simple[]

	// #tag::results[]
	var rows []interface{}
	for results.Next() {
		var row interface{}
		if err := results.Row(&row); err != nil {
			panic(err)
		}
		rows = append(rows, row)
	}

	if err := results.Err(); err != nil {
		panic(err)
	}
	// #end::results[]

	fmt.Println(rows)

	// #tag::metadata[]
	// make sure that results has been iterated (and therefore closed) before calling this.
	metadata, err := results.MetaData()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client context id: %s\n", metadata.ClientContextID)
	fmt.Printf("Elapsed time: %d\n", metadata.Metrics.ElapsedTime)
	fmt.Printf("Execution time: %d\n", metadata.Metrics.ExecutionTime)
	fmt.Printf("Result count: %d\n", metadata.Metrics.ResultCount)
	fmt.Printf("Error count: %d\n", metadata.Metrics.ErrorCount)
	// #end::metadata[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
