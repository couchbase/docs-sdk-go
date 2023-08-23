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
	// #tag::landmarksviewstart[]
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

	viewResult, err := bucket.ViewQuery("landmarks-by-country", "by_country", &gocb.ViewOptions{
		StartKey:        "U",
		Limit:           10,
		Namespace:       gocb.DesignDocumentNamespaceDevelopment,
		ScanConsistency: gocb.ViewScanConsistencyRequestPlus,
	})
	if err != nil {
		panic(err)
	}
	// #end::landmarksviewstart[]

	for viewResult.Next() {
		row := viewResult.Row()
		fmt.Printf("Document ID: %s\n", row.ID)
		var key string
		err = row.Key(&key)
		if err != nil {
			panic(err)
		}

		var landmark interface{}
		err = row.Value(&landmark)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Landmark named %s has value %v\n", key, landmark)
	}

	// always check for errors after iterating
	err = viewResult.Err()
	if err != nil {
		panic(err)
	}

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}

func createView(viewMgr *gocb.ViewIndexManager) {
	designDoc := gocb.DesignDocument{
		Name: "landmarks-by-country",
		Views: map[string]gocb.View{
			"by_country": {
				Map:    "function (doc, meta) { if (doc.type == 'landmark') { emit(doc.country, null); } }",
				Reduce: "",
			},
		},
	}

	err := viewMgr.UpsertDesignDocument(designDoc, gocb.DesignDocumentNamespaceDevelopment, nil)
	if err != nil {
		panic(err)
	}
}
