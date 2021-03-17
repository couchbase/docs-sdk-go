package main

import (
	"github.com/couchbase/gocb/v2"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("travel-sample")

	collection := bucket.DefaultCollection()

	// #tag::rawstring[]
	input := "hello world"
	transcoder := gocb.NewRawStringTranscoder()

	_, err = collection.Upsert("key", input, &gocb.UpsertOptions{
		Transcoder: transcoder,
	})
	if err != nil {
		panic(err)
	}

	getRes, err := collection.Get("key", &gocb.GetOptions{
		Transcoder: transcoder,
	})
	if err != nil {
		panic(err)
	}

	var returned string
	err = getRes.Content(&returned)
	if err != nil {
		panic(err)
	}
	// #end::rawstring[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
