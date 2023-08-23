package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

func operationsWithNewUser(username, password, connString, bucketName string) {
	// tag::operations[]
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	}
	cluster, err := gocb.Connect(connString, opts)
	if err != nil {
		panic(err)
	}
	// For Server versions 6.5 or later you do not need to open a bucket here
	bucket := cluster.Bucket(bucketName)
	collection := bucket.Scope("inventory").Collection("airline")

	err = cluster.QueryIndexes().CreatePrimaryIndex(bucketName, &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
	})
	if err != nil {
		panic(err)
	}

	airline10, err := collection.Get("airline_10", nil)
	if err != nil {
		panic(err)
	}

	var airline interface{}
	err = airline10.Content(&airline)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Airline 10: %v\n", airline)

	airline11 := map[string]interface{}{
		"callsign": "MILE-AIR",
		"iata":     "Q5",
		"id":       11,
		"name":     "40-Mile Air",
		"type":     "airline",
	}
	_, err = collection.Upsert("airline_11", airline11, nil)
	if err != nil {
		panic(err)
	}

	queryRes, err := cluster.Query("SELECT * FROM `travel-sample`.inventory.airline LIMIT 5", nil)
	if err != nil {
		panic(err)
	}

	for queryRes.Next() {
		var queryData interface{}
		err = queryRes.Row(&queryData)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Query row: %v\n", queryData)
	}

	cluster.Close(nil)
	// end::operations[]
}

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	connString := "your-ip"
	cluster, err := gocb.Connect(connString, opts)
	if err != nil {
		panic(err)
	}

	bucketName := "travel-sample"
	// For Server versions 6.5 or later you do not need to open a bucket here
	bucket := cluster.Bucket(bucketName)

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	username := "myusername"
	password := "somethingsecret"

	// tag::upsert[]
	userMgr := cluster.Users()
	user := gocb.User{
		Username:    username,
		DisplayName: "My Displayname",
		Roles: []gocb.Role{
			// Roles required for the reading of data from the bucket
			{
				Name:   "data_reader",
				Bucket: "*",
			},
			{
				Name:   "query_select",
				Bucket: "*",
			},
			// Roles required for the writing of data into the bucket.
			{
				Name:   "data_writer",
				Bucket: bucketName,
			},
			{
				Name:   "query_insert",
				Bucket: bucketName,
			},
			{
				Name:   "query_delete",
				Bucket: bucketName,
			},
			// Role required for the creation of indexes on the bucket.
			{
				Name:   "query_manage_index",
				Bucket: bucketName,
			},
		},
		Password: password,
	}

	err = userMgr.UpsertUser(user, nil)
	if err != nil {
		panic(err)
	}
	// end::upsert[]

	// tag::getall[]
	users, err := userMgr.GetAllUsers(&gocb.GetAllUsersOptions{})
	if err != nil {
		panic(err)
	}

	for _, u := range users {
		fmt.Printf("User's display name is: %s\n", u.DisplayName)
		roles := u.Roles
		for _, r := range roles {
			fmt.Printf("	User has the role %s, applicable to bucket %s\n", r.Name, r.Bucket)
		}
	}
	// end::getall[]

	operationsWithNewUser(username, password, connString, bucketName)

	cluster.Close(nil)
}
