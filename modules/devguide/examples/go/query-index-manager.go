package main

import (
	"errors"
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	// tag::creating-index-mgr[]
	cluster, err := gocb.Connect("localhost", gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	})
	if err != nil {
		panic(err)
	}

	if err = cluster.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}

	queryIndexMgr := cluster.QueryIndexes()
	// end::creating-index-mgr[]

	primaryIndex(queryIndexMgr)
	secondaryIndex(queryIndexMgr)
	deferAndWatchIndex(queryIndexMgr)
	dropPrimaryAndSecondaryIndex(queryIndexMgr)
}

func primaryIndex(queryIndexMgr *gocb.QueryIndexManager) {
	fmt.Println("Example - [primary]")
	// tag::primary[]
	if err := queryIndexMgr.CreatePrimaryIndex("travel-sample",
		&gocb.CreatePrimaryQueryIndexOptions{
			ScopeName:      "tenant_agent_01",
			CollectionName: "users",
			// Set this if you wish to use a custom name
			// CustomName: "custom_name",
			IgnoreIfExists: true,
		},
	); err != nil {
		if errors.Is(err, gocb.ErrIndexExists) {
			fmt.Println("Index already exists")
		} else {
			panic(err)
		}
	}
	// end::primary[]
}

func secondaryIndex(queryIndexMgr *gocb.QueryIndexManager) {
	fmt.Println("\nExample - [secondary]")
	// tag::secondary[]
	if err := queryIndexMgr.CreateIndex("travel-sample", "tenant_agent_01_users_email", []string{"preferred_email"},
		&gocb.CreateQueryIndexOptions{
			ScopeName:      "tenant_agent_01",
			CollectionName: "users",
		},
	); err != nil {
		if errors.Is(err, gocb.ErrIndexExists) {
			fmt.Println("Index already exists")
		} else {
			panic(err)
		}
	}
	// end::secondary[]
}

func deferAndWatchIndex(queryIndexMgr *gocb.QueryIndexManager) {
	fmt.Println("\nExample - [defer-indexes]")
	// tag::defer-indexes[]
	// Create a deferred index
	if err := queryIndexMgr.CreateIndex("travel-sample", "tenant_agent_01_users_phone", []string{"preferred_phone"},
		&gocb.CreateQueryIndexOptions{
			ScopeName:      "tenant_agent_01",
			CollectionName: "users",
			Deferred:       true,
		},
	); err != nil {
		if errors.Is(err, gocb.ErrIndexExists) {
			fmt.Println("Index already exists")
		} else {
			panic(err)
		}
	}

	// Build any deferred indexes within `travel-sample`.tenant_agent_01.users
	indexesToBuild, err := queryIndexMgr.BuildDeferredIndexes("travel-sample",
		&gocb.BuildDeferredQueryIndexOptions{
			ScopeName:      "tenant_agent_01",
			CollectionName: "users",
		},
	)
	if err != nil {
		panic(err)
	}

	// Wait for indexes to come online
	if err = queryIndexMgr.WatchIndexes("travel-sample", indexesToBuild, time.Duration(30*time.Second),
		&gocb.WatchQueryIndexOptions{
			ScopeName:      "tenant_agent_01",
			CollectionName: "users",
		},
	); err != nil {
		panic(err)
	}
	// end::defer-indexes[]
}

func dropPrimaryAndSecondaryIndex(queryIndexMgr *gocb.QueryIndexManager) {
	fmt.Println("\nExample - [drop-primary-or-secondary-index]")
	// tag::drop-primary-or-secondary-index[]
	// Drop a primary index
	if err := queryIndexMgr.DropPrimaryIndex("travel-sample",
		&gocb.DropPrimaryQueryIndexOptions{
			ScopeName:      "tenant_agent_01",
			CollectionName: "users",
		},
	); err != nil {
		panic(err)
	}

	// Drop a secondary index
	if err := queryIndexMgr.DropIndex("travel-sample", "tenant_agent_01_users_email",
		&gocb.DropQueryIndexOptions{
			ScopeName:      "tenant_agent_01",
			CollectionName: "users",
		},
	); err != nil {
		panic(err)
	}
	// end::drop-primary-or-secondary-index[]
}
