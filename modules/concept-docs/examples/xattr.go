package main

import "github.com/couchbase/gocb/v2"

func main() {
	// Connect to Cluster
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	cluster, err := gocb.Connect("your-ip", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("travel-sample")
	collection := bucket.Scope("inventory").Collection("airport")

	// tag::xattr[]
	_, err = collection.LookupIn("airport_1254", []gocb.LookupInSpec{
		gocb.GetSpec("$document.exptime", &gocb.GetSpecOptions{IsXattr: true}),
	}, nil)
	if err != nil {
		panic(err)
	}
	// end::xattr[]
}
