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

	// tag::connect[]
	// Update this to your cluster details
	endpoint := "cb.<your-endpoint>.cloud.couchbase.com"
	bucketName := "travel-sample"
	username := "username"
	password := "Password123!"

	// Initialize the Connection
	cluster, err := gocb.Connect("couchbases://"+endpoint, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})
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

	// tag::document[]
	type User struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Interests []string `json:"interests"`
	}
	// end::document[]

	// tag::upsert[]
	// Create and store a Document
	_, err = col.Upsert("u:kingarthur",
		User{
			Name:      "Arthur",
			Email:     "kingarthur@couchbase.com",
			Interests: []string{"Holy Grail", "African Swallows"},
		}, nil)
	if err != nil {
		log.Fatal(err)
	}
	// end::upsert[]

	// tag::get[]
	// Get the document back
	getResult, err := col.Get("u:kingarthur", nil)
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
	queryResult, err := cluster.Query(
		fmt.Sprintf("SELECT name FROM `%s` WHERE $1 IN interests", bucketName),
		&gocb.QueryOptions{PositionalParameters: []interface{}{"African Swallows"}},
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
