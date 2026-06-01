package main

import (
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/sirupsen/logrus"

	"student-record/internal"
)

func main() {
	connectionString := "<<connection-string>>" // Replace this with Connection String
	username := "<<username>>"                  // Replace this with username from cluster access credentials
	password := "<<password>>"                  // Replace this with password from cluster access credentials

	// Setup info level logging.
	gocb.SetLogger(internal.NewLogger(logrus.InfoLevel))

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

	// The code here is similar to creating a student record, but it writes to a different collection.
	courseRecords := scope.Collection("course-record-collection")

	addCourse(courseRecords, "ART-HISTORY-000001", "art history", "fine art", 100)
	addCourse(courseRecords, "FINE-ART-000002", "fine art", "fine art", 50)
	addCourse(courseRecords, "GRAPHIC-DESIGN-000003", "graphic design", "media and communication", 200)

	err = cluster.Close(nil)
	if err != nil {
		log.Fatal(err)
	}
}

func addCourse(collection *gocb.Collection, id, name, faculty string, creditPoints int) {
	course := map[string]interface{}{
		"course-name":   name,
		"faculty":       faculty,
		"credit-points": creditPoints,
	}

	_, err := collection.Upsert(id, course, nil)
	if err != nil {
		log.Fatal(err)
	}
}
