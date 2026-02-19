package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/couchbase/gocb/v2"
)

var bucketName = "travel-sample"

func main() {
	// Connect to Cluster
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

	bucket := cluster.Bucket("travel-sample")
	collection := bucket.Scope("inventory").Collection("airport")

	if err := bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}

	{
	    fmt.Println("Example - [prepared-statement]")
		// tag::prepared-statement[]
		query := "SELECT count(*) FROM `travel-sample`.inventory.airport where country = $1;"
		rows, err := cluster.Query(query, &gocb.QueryOptions{
			Adhoc:                false,
			PositionalParameters: []interface{}{"France"},
		})
		if err != nil {
			panic(err)

		}

		for rows.Next() {
			// do something
		}
		if err := rows.Err(); err != nil {
			panic(err)
		}
		// end::prepared-statement[]
	}

	{
	    fmt.Println("Example - [create-index]")
		// tag::create-index[]
		mgr := cluster.QueryIndexes()
		if err := mgr.CreatePrimaryIndex(bucketName, nil); err != nil {
			if errors.Is(err, gocb.ErrIndexExists) {
				fmt.Println("Index already exists")
			} else {
				panic(err)
			}
		}

		if err := mgr.CreateIndex(bucketName, "ix_name", []string{"name"}, nil); err != nil {
			if errors.Is(err, gocb.ErrIndexExists) {
				fmt.Println("Index already exists")
			} else {
				panic(err)
			}
		}

		if err := mgr.CreateIndex(bucketName, "ix_email", []string{"email"}, nil); err != nil {
			if errors.Is(err, gocb.ErrIndexExists) {
				fmt.Println("Index already exists")
			} else {
				panic(err)
			}
		}
		// end::create-index[]
	}

	{
	    fmt.Println("Example - [deferred-index]")
		// tag::deferred-index[]
		mgr := cluster.QueryIndexes()
		if err := mgr.CreatePrimaryIndex(bucketName,
			&gocb.CreatePrimaryQueryIndexOptions{Deferred: true},
		); err != nil {
			if errors.Is(err, gocb.ErrIndexExists) {
				fmt.Println("Index already exists")
			} else {
				panic(err)
			}
		}

		if err := mgr.CreateIndex(bucketName, "ix_name", []string{"name"},
			&gocb.CreateQueryIndexOptions{Deferred: true},
		); err != nil {
			if errors.Is(err, gocb.ErrIndexExists) {
				fmt.Println("Index already exists")
			} else {
				panic(err)
			}
		}

		if err = mgr.CreateIndex(bucketName, "ix_email", []string{"email"},
			&gocb.CreateQueryIndexOptions{Deferred: true},
		); err != nil {
			if errors.Is(err, gocb.ErrIndexExists) {
				fmt.Println("Index already exists")
			} else {
				panic(err)
			}
		}

		indexesToBuild, err := mgr.BuildDeferredIndexes(bucketName, nil)
		if err != nil {
			panic(err)
		}
		err = mgr.WatchIndexes(bucketName, indexesToBuild, time.Duration(2*time.Second), nil)
		if err != nil {
			panic(err)
		}
		// end::deferred-index[]
	}

	{
	    fmt.Println("Example - [index-consistency]")
		// tag::index-consistency[]
		random := rand.Intn(10000000)
		user := struct {
			Name   string `json:"name"`
			Email  string `json:"email"`
			Random int    `json:"random"`
		}{Name: "Brass Doorknob", Email: "brass.doorknob@juno.com", Random: random}

		_, err := collection.Upsert(fmt.Sprintf("user:%d", random), user, nil)
		if err != nil {
			panic(err)
		}

		_, err = cluster.Query(
			"SELECT name, email, random, META().id FROM `travel-sample`.inventory.airport WHERE $1 IN name",
			&gocb.QueryOptions{
				PositionalParameters: []interface{}{"Brass"},
			},
		)
		if err != nil {
			panic(err)
		}
		// end::index-consistency[]
	}

	{
	    fmt.Println("Example - [index-consistency-request-plus]")
		// tag::index-consistency-request-plus[]
		_, err := cluster.Query(
			"SELECT name, email, random, META().id FROM `travel-sample`.inventory.airport WHERE $1 IN name",
			&gocb.QueryOptions{
				PositionalParameters: []interface{}{"Brass"},
				ScanConsistency:      gocb.QueryScanConsistencyRequestPlus,
			},
		)
		if err != nil {
			panic(err)
		}
		// end::index-consistency-request-plus[]
	}
}
