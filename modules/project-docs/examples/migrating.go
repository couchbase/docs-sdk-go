package main

import (
	"crypto/tls"
	"errors"
	"log"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func basic() {
	// #tag::basicconnecting[]
	opts := gocb.ClusterOptions{
		Username: "Administrator",
		Password: "password",
	}
	cluster, err := gocb.Connect("10.112.193.101", opts)
	if err != nil {
		panic(err)
	}
	// #end::basicconnecting[]

	// #tag::getcollection[]
	bucket := cluster.Bucket("travel-sample")
	collection := bucket.Scope("inventory").Collection("airport")
	// #end::getcollection[]

	// #tag::getresult[]
	getResult, err := collection.Get("key", &gocb.GetOptions{
		Timeout: 2 * time.Second,
	})
	// #end::getresult[]
	// #tag::handleerror[]
	if errors.Is(err, gocb.ErrDocumentNotFound) {
		// handle your error
	}
	// #end::handleerror[]
	// #tag::handleerrorext[]
	if errors.Is(err, gocb.ErrDocumentNotFound) {
		var kverr gocb.KeyValueError
		if errors.As(err, &kverr) {
			log.Printf("Error Context: %+v\n", kverr)
		}
	}
	// #end::handleextendederror[]
	log.Printf("Get Result: %+v\n", getResult)

	// #tag::queryresult[]
	queryResult, err := cluster.Query("select 1=1", &gocb.QueryOptions{
		Timeout: 3 * time.Second,
	})
	// #end::queryresult[]
	log.Printf("Query Result: %+v\n", queryResult)

	cluster.Close(&gocb.ClusterCloseOptions{})
}

func passauthenticator() {
	// #tag::passauthenticator[]
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	cluster, err := gocb.Connect("10.112.193.101", opts)
	if err != nil {
		panic(err)
	}
	// #end::passauthenticator[]

	cluster.Close(&gocb.ClusterCloseOptions{})
}

func certauthenticator() {
	// #tag::certauthenticator[]
	cert, err := tls.LoadX509KeyPair("mycert.pem", "mykey.pem")
	if err != nil {
		panic(err)
	}

	opts := gocb.ClusterOptions{
		Authenticator: gocb.CertificateAuthenticator{
			ClientCertificate: &cert,
		},
	}
	cluster, err := gocb.Connect("10.112.193.101", opts)
	if err != nil {
		panic(err)
	}
	// #end::certauthenticator[]

	cluster.Close(&gocb.ClusterCloseOptions{})
}

func main() {
	basic()
	passauthenticator()
	certauthenticator()
}
