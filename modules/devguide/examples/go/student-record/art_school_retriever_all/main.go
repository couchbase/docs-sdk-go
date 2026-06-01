package main

import (
	"fmt"
	"log"

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

	retrieveCourses(cluster)

	err = cluster.Close(nil)
	if err != nil {
		log.Fatal(err)
	}
}

func retrieveCourses(cluster *gocb.Cluster) {
	queryResult, err := cluster.Query(
		"SELECT crc.* FROM `student-bucket`.`art-school-scope`.`course-record-collection` crc",
		&gocb.QueryOptions{},
	)
	if err != nil {
		log.Fatal(err)
	}

	for queryResult.Next() {
		var row interface{}
		err := queryResult.Row(&row)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Found row:", row)
	}

	if err := queryResult.Err(); err != nil {
		log.Fatal(err)
	}
}
