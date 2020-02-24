package main

import (
	"github.com/couchbase/gocb/v2"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}

	// #tag::creatingbucketmgr[]
	cluster, err := gocb.Connect("10.112.193.101", opts)
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
			RAMQuotaMB:           1024,
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
	settings, err := bucketMgr.GetBucket("test", nil)
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
	err := bucketMgr.DropBucket("test", nil)
	if err != nil {
		panic(err)
	}
	// #end::removeBucket[]
}

func flushBucket(bucketMgr *gocb.BucketManager) {
	// #tag::flushBucket[]
	err := bucketMgr.FlushBucket("test", nil)
	if err != nil {
		panic(err)
	}
	// #end::flushBucket[]
}
