package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// Uncomment following line to enable logging
	// gocb.SetLogger(gocb.VerboseStdioLogger())
	endpoint := "cb.<your endpoint address>.dp.cloud.couchbase.com"
	bucketName := "couchbasecloudbucket"
	username := "user"
	password := "password"

	// Initialize the Connection
	cluster, err := gocb.Connect("couchbases://"+endpoint+"?ssl=no_verify", gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	col := bucket.DefaultCollection()

	// Create a N1QL Primary Index (but ignore if it exists)
	cluster.QueryIndexes().CreatePrimaryIndex(bucketName, &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
	})

	type User struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Interests []string `json:"interests"`
	}

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
			panic(err)
		}
		fmt.Println(result)
	}

	if err := queryResult.Err(); err != nil {
		panic(err)
	}

}
