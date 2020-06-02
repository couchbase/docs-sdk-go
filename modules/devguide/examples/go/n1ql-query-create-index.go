package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
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

	idxMgr := cluster.QueryIndexes()

	// Create a primary index
	err = idxMgr.CreatePrimaryIndex("travel-sample", &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
	})
	if err != nil {
		panic(err)
	}

	// Create a deferred named index
	err = idxMgr.CreateIndex("travel-sample", "my-index", []string{"name"}, &gocb.CreateQueryIndexOptions{
		IgnoreIfExists: true,
		Deferred:       true,
	})
	if err != nil {
		panic(err)
	}

	// Build deferred indexes
	built, err := idxMgr.BuildDeferredIndexes("travel-sample", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Building indexes: %v", built)

	// Drop the primary index
	err = idxMgr.DropPrimaryIndex("travel-sample", nil)
	if err != nil {
		panic(err)
	}

	// Drop the named index
	err = idxMgr.DropIndex("travel-sample", "my-index", nil)
	if err != nil {
		panic(err)
	}

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
