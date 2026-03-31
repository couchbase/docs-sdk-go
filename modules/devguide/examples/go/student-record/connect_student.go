package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	connectionString := "<<connection-string>>" // Replace this with Connection String
	username := "<<username>>"                  // Replace this with username from cluster access credentials
	password := "<<password>>"                  // Replace this with password from cluster access credentials

	// Setup info level logging.
	gocb.SetLogger(NewLogger(logrus.InfoLevel))

	// Connecting to the cluster
	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	}

	// Use the pre-configured profile below to avoid latency issues with your connection.
	if err := options.ApplyProfile(gocb.ClusterConfigProfileWanDevelopment); err != nil {
		log.Fatal(err)
	}

	cluster, err := gocb.Connect("couchbases://"+connectionString, options)
	if err != nil {
		log.Fatal(err)
	}

	// The `cluster.Bucket` retrieves the bucket you set up for the student cluster.
	bucket := cluster.Bucket("student-bucket")

	// Forces the application to wait until the bucket is ready.
	err = bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	// The `bucket.Scope` retrieves the `art-school-scope` from the bucket.
	scope := bucket.Scope("art-school-scope")

	// The `scope.Collection` retrieves the student collection from the scope.
	studentRecords := scope.Collection("student-record-collection")

	// A check to make sure the collection is connected and retrieved when you run the application.
	fmt.Println("The name of this collection is", studentRecords.Name())

	// Like with all database systems, it's good practice to disconnect from the Couchbase cluster after you have finished working with it.
	err = cluster.Close(nil)
	if err != nil {
		log.Fatal(err)
	}
}
