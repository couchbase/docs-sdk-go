package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
	connectionString := "<<connection-string>>" // Replace this with Connection String
	username := "<<username>>"                  // Replace this with username from cluster access credentials
	password := "<<password>>"                  // Replace this with password from cluster access credentials

	// Setup info level logging.
	gocb.SetLogger(NewLogger(logrus.InfoLevel))

	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	}

	if err := options.ApplyProfile(gocb.ClusterConfigProfileWanDevelopment); err != nil {
		log.Fatal(err)
	}

	cluster, err := gocb.Connect("couchbases://"+connectionString, options)
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket("student-bucket")
	err = bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	scope := bucket.Scope("art-school-scope")
	studentRecords := scope.Collection("student-record-collection")

	// Create and populate the student record.
	hilary := map[string]interface{}{
		"name":          "Hilary Smith",
		"date-of-birth": "1980-12-21",
	}

	// The `Upsert` function inserts or updates documents in a collection.
	// The first parameter is a unique ID for the document, similar to a primary key used in a relational database system.
	// If the `Upsert` call finds a document with a matching ID in the collection, it updates the document.
	// If there is no matching ID, it creates a new document.
	_, err = studentRecords.Upsert("000001", hilary, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = cluster.Close(nil)
	if err != nil {
		log.Fatal(err)
	}
}
