package main

import (
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

	// tag::buckets_and_clusters_1[]
	bucketMgr := cluster.Buckets()
	createBucketSettings := gocb.CreateBucketSettings{
		BucketSettings: gocb.BucketSettings{
			Name:                 "myBucket",
			RAMQuotaMB:           150,
			BucketType:           gocb.CouchbaseBucketType,
		},
	}
	if err := bucketMgr.CreateBucket(createBucketSettings, &gocb.CreateBucketOptions{}); err != nil {
		panic(err)
	}
	// end::buckets_and_clusters_1[]
}
