package main

import (
	"fmt"

	gocb "github.com/couchbase/gocb/v2"
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

	// get a bucket reference
	bucket := cluster.Bucket("travel-sample")

	// #tag::landmarksview[]
	landmarksResult, err := bucket.ViewQuery("landmarks", "by_name", &gocb.ViewOptions{
		Key:       "<landmark-name>",
		Namespace: gocb.DesignDocumentNamespaceDevelopment,
	})
	if err != nil {
		panic(err)
	}
	// #end::landmarksview[]

	// #tag::results[]
	for landmarksResult.Next() {
		landmarkRow := landmarksResult.Row()
		fmt.Printf("Document ID: %s\n", landmarkRow.ID)
		var key string
		err = landmarkRow.Key(&key)
		if err != nil {
			panic(err)
		}

		var landmark interface{}
		err = landmarkRow.Value(&landmark)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Landmark named %s has value %v\n", key, landmark)
	}

	// always check for errors after iterating
	err = landmarksResult.Err()
	if err != nil {
		panic(err)
	}
	// #end::results[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
