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

	// Retrieves the student bucket you set up.
	bucket := cluster.Bucket("student-bucket")

	// Forces the application to wait until the bucket is ready.
	err = bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieves the `art-school-scope` collection from the scope.
	scope := bucket.Scope("art-school-scope")
	studentRecords := scope.Collection("student-record-collection")

	// Retrieves Hilary's student record, the `graphic design` course record, and the `art history` course record.
	// Each method uses a SQL++ call to retrieve a single record from each collection.
	hilary := retrieveStudent(cluster, "Hilary Smith")
	graphicDesign := retrieveCourse(cluster, "graphic design")
	artHistory := retrieveCourse(cluster, "art history")

	// Couchbase does not have a native date type, so the common practice is to store dates as strings.
	currentDate := time.Now().Format("2006-01-02")

	// Stores the `enrollments` inside the student record as an array.
	enrollments := []map[string]interface{}{
		{
			"course-id":     graphicDesign["id"],
			"date-enrolled": currentDate,
		},
		{
			"course-id":     artHistory["id"],
			"date-enrolled": currentDate,
		},
	}

	// Adds the `enrollments` array to Hilary's student record.
	hilary["enrollments"] = enrollments

	// Commits the changes to the collection.
	// The `Upsert` function call takes the key of the record you want to insert or update and the record itself as parameters.
	// If the `Upsert` call finds a document with a matching ID in the collection, it updates the document.
	// If there is no matching ID, it creates a new document.
	studentID, ok := hilary["id"].(string)
	if !ok {
		log.Fatal("student record missing id field")
	}
	delete(hilary, "id")
	_, err = studentRecords.Upsert(studentID, hilary, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = cluster.Close(nil)
	if err != nil {
		log.Fatal(err)
	}
}

func retrieveStudent(cluster *gocb.Cluster, name string) map[string]interface{} {
	queryResult, err := cluster.Query(
		"SELECT META().id, src.* FROM `student-bucket`.`art-school-scope`.`student-record-collection` src WHERE src.`name` = $name",
		&gocb.QueryOptions{
			NamedParameters: map[string]interface{}{
				"name": name,
			},
			ScanConsistency: gocb.QueryScanConsistencyRequestPlus,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	var row map[string]interface{}
	if queryResult.Next() {
		err = queryResult.Row(&row)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := queryResult.Err(); err != nil {
		log.Fatal(err)
	}

	return row
}

func retrieveCourse(cluster *gocb.Cluster, course string) map[string]interface{} {
	queryResult, err := cluster.Query(
		"SELECT META().id, crc.* FROM `student-bucket`.`art-school-scope`.`course-record-collection` crc WHERE crc.`course-name` = $courseName",
		&gocb.QueryOptions{
			NamedParameters: map[string]interface{}{
				"courseName": course,
			},
			ScanConsistency: gocb.QueryScanConsistencyRequestPlus,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	var row map[string]interface{}
	if queryResult.Next() {
		err = queryResult.Row(&row)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := queryResult.Err(); err != nil {
		log.Fatal(err)
	}

	return row
}
