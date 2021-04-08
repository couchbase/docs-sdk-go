package main

import (
	"errors"
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func GetCollections(username string, password string) *gocb.CollectionManager {
	fmt.Println("create-collection-manager")

	// tag::create-collection-manager[]
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("travel-sample")
	collections := bucket.Collections()
	// end::create-collection-manager[]

	return collections
}

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

	// For Server versions 6.5 or later you do not need to open a bucket here
	b := cluster.Bucket("travel-sample")

	// We wait until the bucket is definitely connected and setup.
	// For Server versions 6.5 or later if we hadn't opened a bucket then we could use cluster.WaitUntilReady here.
	err = b.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	users := cluster.Users()

	fmt.Println("bucketAdmin")
	// tag::bucketAdmin[]
	{
		user := gocb.User{
			Username:    "bucketAdmin",
			DisplayName: "Bucket Admin [travel-sample]",
			Password:    "password",
			Roles: []gocb.Role{
				{
					Name:   "bucket_admin",
					Bucket: "travel-sample",
				}}}

		err = users.UpsertUser(user, nil)
		if err != nil {
			panic(err)
		}
	}
	// end::bucketAdmin[]

	{
		collections := GetCollections("bucketAdmin", "password")

		fmt.Println("create-scope")
		// tag::create-scope[]
		err = collections.CreateScope("example-scope", nil)
		if err != nil {
			if errors.Is(err, gocb.ErrScopeExists) {
				fmt.Println("Scope already exists")
			} else {
				panic(err)
			}
		}
		// end::create-scope[]
	}

	fmt.Println("scopeAdmin")
	// tag::scopeAdmin[]
	{
		user := gocb.User{
			Username:    "scopeAdmin",
			DisplayName: "Manage Collections in Scope [travel-sample:*]",
			Password:    "password",
			Roles: []gocb.Role{
				{
					Name:   "scope_admin",
					Bucket: "travel-sample",
					Scope:  "example-scope",
				},
				{
					Name:   "data_reader",
					Bucket: "travel-sample",
				},
			}}

		err = users.UpsertUser(user, nil)
		if err != nil {
			panic(err)
		}
	}
	// end::scopeAdmin[]

	{
		collections := GetCollections("scopeAdmin", "password")

		{
			fmt.Println("create-collection")
			// tag::create-collection[]
			collection := gocb.CollectionSpec{
				Name:      "example-collection",
				ScopeName: "example-scope",
			}

			err = collections.CreateCollection(collection, nil)
			if err != nil {
				if errors.Is(err, gocb.ErrCollectionExists) {
					fmt.Println("Collection already exists")
				} else {
					panic(err)
				}
			}
			// end::create-collection[]

			fmt.Println("drop-collection")
			// tag::drop-collection[]
			err = collections.DropCollection(collection, nil)
			if err != nil {
				panic(err)
			}
			// end::drop-collection[]
		}

		{
			fmt.Println("drop-scope")
			// tag::drop-scope[]
			err = collections.DropScope("example-scope", nil)
			if err != nil {
				panic(err)
			}
			// end::drop-scope[]
		}

	}

	cluster.Close(nil)
}
