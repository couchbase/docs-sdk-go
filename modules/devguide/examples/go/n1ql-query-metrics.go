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

	// #tag::metrics[]
	query := "SELECT x.* FROM `travel-sample`.inventory.airport x LIMIT 10;"
	rows, err := cluster.Query(query, &gocb.QueryOptions{
		Metrics: true,
		Adhoc:   true,
	})
	// check query was successful
	if err != nil {
		panic(err)
	}

	// iterate over rows
	for rows.Next() {
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
	// #end::metrics[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
