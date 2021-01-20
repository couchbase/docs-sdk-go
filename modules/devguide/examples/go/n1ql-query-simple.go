package main

import (
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
	cluster, err := gocb.Connect("localhost", opts)
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

	// #tag::results[]
	query := "SELECT x.* FROM `travel-sample` x LIMIT 10;"
	rows, err := cluster.Query(query, &gocb.QueryOptions{})
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

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
