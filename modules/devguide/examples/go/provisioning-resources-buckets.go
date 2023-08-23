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

	// #tag::creatingbucketmgr[]
	cluster, err := gocb.Connect("your-ip opts)
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

	bucketMgr := cluster.Buckets()

	// #end::creatingbucketmgr[]

	createBucket(bucketMgr)
	updateBucket(bucketMgr)
	flushBucket(bucketMgr)
	removeBucket(bucketMgr)
}

func createBucket(bucketMgr *gocb.BucketManager) {
	// #tag::createBucket[]
	err := bucketMgr.CreateBucket(gocb.CreateBucketSettings{
		BucketSettings: gocb.BucketSettings{
			Name:                 "hello",
			FlushEnabled:         false,
			ReplicaIndexDisabled: true,
			RAMQuotaMB:           150,
			NumReplicas:          1,
			BucketType:           gocb.CouchbaseBucketType,
		},
		ConflictResolutionType: gocb.ConflictResolutionTypeSequenceNumber,
	}, nil)
	if err != nil {
		panic(err)
	}
	// #end::createBucket[]
}

func updateBucket(bucketMgr *gocb.BucketManager) {
	// #tag::updateBucket[]
	settings, err := bucketMgr.GetBucket("hello", nil)
	if err != nil {
		panic(err)
	}

	settings.FlushEnabled = true
	err = bucketMgr.UpdateBucket(*settings, nil)
	if err != nil {
		panic(err)
	}
	// #end::updateBucket[]
}

func removeBucket(bucketMgr *gocb.BucketManager) {
	// #tag::removeBucket[]
	err := bucketMgr.DropBucket("hello", nil)
	if err != nil {
		panic(err)
	}
	// #end::removeBucket[]
}

func flushBucket(bucketMgr *gocb.BucketManager) {
	// allow some time before flushing the bucket as example code
	// runs fairly quickly.
	time.Sleep(5 * time.Second)

	// #tag::flushBucket[]
	err := bucketMgr.FlushBucket("hello", nil)
	if err != nil {
		panic(err)
	}
	// #end::flushBucket[]
}
