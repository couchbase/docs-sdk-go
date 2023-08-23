package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// Connect to Cluster
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	cluster, err := gocb.Connect("your-ip opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("travel-sample")

	// We wait until the bucket is definitely connected and setup.
	// For Server versions 6.5 or later if we hadn't opened a bucket then we could use cluster.WaitUntilReady here.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	// tag::namedcollectionupsert[]
	agentScope := bucket.Scope("tenant_agent_00")
	usersCollection := agentScope.Collection("users")

	type userDoc struct {
		Name           string `json:"name"`
		PreferredEmail string `json:"preferred_email"`
	}
	document := userDoc{Name: "John Doe", PreferredEmail: "johndoe111@test123.test"}

	result, err := usersCollection.Upsert("user-key", &document, &gocb.UpsertOptions{})
	if err != nil {
		panic(err)
	}
	// end::namedcollectionupsert[]
	fmt.Println(result.Cas())

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
