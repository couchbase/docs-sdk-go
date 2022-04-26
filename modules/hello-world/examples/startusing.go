package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// Uncomment following line to enable logging
	// gocb.SetLogger(gocb.VerboseStdioLogger())

	// tag::connect[]
	// Update this to your cluster details
	bucketName := "travel-sample"
	username := "Administrator"
	password := "password"

	p, err := ioutil.ReadFile("path/to/ca.pem")
	if err != nil {
		log.Fatal(err)
	}

	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(p)

	// Initialize the Connection
	cluster, err := gocb.Connect("couchbases://127.0.0.1", gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
		SecurityConfig: gocb.SecurityConfig{
			TLSRootCAs: roots,
			// WARNING: DO not set this to true in production, only use this for testing!
			// TLSSkipVerify: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	// end::connect[]

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get a reference to the default collection, required for older Couchbase server versions
	// col := bucket.DefaultCollection()

	col := bucket.Scope("tenant_agent_00").Collection("users")

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
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	if err := queryResult.Err(); err != nil {
		log.Fatal(err)
	}
}
