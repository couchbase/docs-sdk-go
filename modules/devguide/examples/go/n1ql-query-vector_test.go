package main

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
)

func Example_n1qlQueryVector() {
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

	{
		// #tag::vector-hyperscale-index[]
		query := "SELECT d.id, d.question, d.wanted_similar_color_from_search, " +
			"ARRAY_CONCAT( " +
			"d.couchbase_search_query.knn[0].vector[0:4], " +
			"['...'] " +
			") AS vector " +
			"FROM `vector-sample`.`color`.`rgb-questions` AS d " +
			"WHERE d.id = '#87CEEB';"

		rows, err := cluster.Query(query, &gocb.QueryOptions{Metrics: true})

		// check query was successful
		if err != nil {
			panic(err)
		}

		// iterate over rows
		for rows.Next() {
			var r interface{}
			err := rows.Row(&r)
			if err != nil {
				panic(err)
			}
			fmt.Println(r)
		}

		// always check for errors after iterating
		err = rows.Err()
		if err != nil {
			panic(err)
		}
		// #end::vector-hyperscale-index[]
	}

	{
		// #tag::vector-hyperscale-index-parameterized[]
		query := "SELECT d.id, d.question, d.wanted_similar_color_from_search, " +
			"ARRAY_CONCAT( " +
			"d.couchbase_search_query.knn[0].vector[0:4], " +
			"['...'] " +
			") AS vector " +
			"FROM `vector-sample`.`color`.`rgb-questions` AS d " +
			"WHERE d.id = $id;"

		rows, err := cluster.Query(query, &gocb.QueryOptions{
			NamedParameters: map[string]interface{}{
				"id": "#87CEEB",
			},
		})

		// check query was successful
		if err != nil {
			panic(err)
		}

		// iterate over rows
		for rows.Next() {
			var r interface{}
			err := rows.Row(&r)
			if err != nil {
				panic(err)
			}
			fmt.Println(r)
		}

		// always check for errors after iterating
		err = rows.Err()
		if err != nil {
			panic(err)
		}
		// #end::vector-hyperscale-index-parameterized[]
	}

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
