// tag::imports[]
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

// end::imports[]

func main() {
	// Uncomment following line to enable logging
	// gocb.SetLogger(gocb.VerboseStdioLogger())

	// tag::connect-info[]
	// Update this to your cluster details
	connectionString := "cb.<your-endpoint>.cloud.couchbase.com"
	bucketName := "travel-sample"
	username := "username"
	password := "Password!123"
	// end::connect-info[]

	// tag::connect[]
	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	}

	// Sets a pre-configured profile called "wan-development" to help avoid latency issues
	// when accessing Capella from a different Wide Area Network
	// or Availability Zone (e.g. your laptop).
	if err := options.ApplyProfile(gocb.ClusterConfigProfileWanDevelopment); err != nil {
		log.Fatal(err)
	}

	// Initialize the Connection
	cluster, err := gocb.Connect("couchbases://"+connectionString, options)
	if err != nil {
		log.Fatal(err)
	}
	// end::connect[]

	// tag::bucket[]
	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}
	// end::bucket[]

	// tag::collection[]
	// Get a reference to the default collection, required for older Couchbase server versions
	// col := bucket.DefaultCollection()

	col := bucket.Scope("tenant_agent_00").Collection("users")
	// end::collection[]

	// tag::upsert[]
	// Create and store a Document
	type User struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Interests []string `json:"interests"`
	}

	_, err = col.Upsert("u:jade",
		User{
			Name:      "Jade",
			Email:     "jade@test-email.com",
			Interests: []string{"Swimming", "Rowing"},
		}, nil)
	if err != nil {
		log.Fatal(err)
	}
	// end::upsert[]

	// tag::get[]
	// Get the document back
	getResult, err := col.Get("u:jade", nil)
	if err != nil {
		log.Fatal(err)
	}

	var inUser User
	err = getResult.Content(&inUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %v\n", inUser)
	// end::get[]

	// tag::query[]
	// Perform a N1QL Query
	inventoryScope := bucket.Scope("inventory")
	queryResult, err := inventoryScope.Query(
		fmt.Sprintf("SELECT * FROM airline WHERE id=10"),
		&gocb.QueryOptions{},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print each found Row
	for queryResult.Next() {
		var result interface{}
		err := queryResult.Row(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	if err := queryResult.Err(); err != nil {
		log.Fatal(err)
	}
	// end::query[]
}
