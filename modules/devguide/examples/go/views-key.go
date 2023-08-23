package main

import (
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
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

	// get a bucket reference
	bucket := cluster.Bucket("travel-sample")

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	// creates the required view for the query
	viewMgr := bucket.ViewIndexes()
	createView(viewMgr)

	// #tag::landmarksview[]
	landmarksResult, err := bucket.ViewQuery("landmarks-by-name", "by_name", &gocb.ViewOptions{
		Key:             "Circle Bar",
		Namespace:       gocb.DesignDocumentNamespaceDevelopment,
		ScanConsistency: gocb.ViewScanConsistencyRequestPlus,
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

func createView(viewMgr *gocb.ViewIndexManager) {
	designDoc := gocb.DesignDocument{
		Name: "landmarks-by-name",
		Views: map[string]gocb.View{
			"by_name": {
				Map:    "function (doc, meta) { if (doc.type == 'landmark') { emit(doc.name, null); } }",
				Reduce: "",
			},
		},
	}

	err := viewMgr.UpsertDesignDocument(designDoc, gocb.DesignDocumentNamespaceDevelopment, nil)
	if err != nil {
		panic(err)
	}
}
