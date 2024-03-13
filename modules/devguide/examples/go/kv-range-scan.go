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
	cluster, err := gocb.Connect("couchbase://localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("default")
	collection := bucket.DefaultCollection()

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}
	{
		// #tag::range_scan_all_documents[]
		results, err := collection.Scan(gocb.RangeScan{}, nil) // <1>
		if err != nil {
			panic(err)
		}
		for {
			item := results.Next()
			if item == nil {
				break
			}
			var content interface{}
			err = item.Content(&content)
			if err != nil {
				panic(err)
			}
			fmt.Printf("ID = %s, \tContent = %s\n", item.ID(), content)
		}
		// Always check for errors after iterating
		err = results.Err()
		if err != nil {
			panic(err)
		}
		// #end::range_scan_all_documents[]
	}
	{
		// #tag::range_scan_all_document_ids[]
		results, err := collection.Scan(gocb.RangeScan{}, &gocb.ScanOptions{IDsOnly: true})
		if err != nil {
			panic(err)
		}
		for {
			item := results.Next()
			if item == nil {
				break
			}
			fmt.Printf("ID = %s\n", item.ID())
		}
		// Always check for errors after iterating
		err = results.Err()
		if err != nil {
			panic(err)
		}
		// #end::range_scan_all_document_ids[]
	}
	{
		// #tag::range_scan_prefix[]
		results, err := collection.Scan(gocb.NewRangeScanForPrefix("alice::"), nil) // <1>
		if err != nil {
			panic(err)
		}
		for {
			item := results.Next()
			if item == nil {
				break
			}
			var content interface{}
			err = item.Content(&content)
			if err != nil {
				panic(err)
			}
			fmt.Printf("ID = %s, \tContent = %s\n", item.ID(), content)
		}
		// Always check for errors after iterating
		err = results.Err()
		if err != nil {
			panic(err)
		}
		// #end::range_scan_prefix[]
	}
	{
		// #tag::range_scan_sample[]
		results, err := collection.Scan(gocb.SamplingScan{Limit: 100}, nil)
		if err != nil {
			panic(err)
		}
		for {
			item := results.Next()
			if item == nil {
				break
			}
			var content interface{}
			err = item.Content(&content)
			if err != nil {
				panic(err)
			}
			fmt.Printf("ID = %s, \tContent = %s\n", item.ID(), content)
		}
		// Always check for errors after iterating
		err = results.Err()
		if err != nil {
			panic(err)
		}
		// #end::range_scan_sample[]
	}

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
