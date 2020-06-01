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
	rows, err := cluster.AnalyticsQuery("select airportname, country from airports where country = 'France' LIMIT 1;", nil)
	// check query was successful
	if err != nil {
		panic(err)
	}

	var result interface{}
	err = rows.One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// #end::one[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
