package main

import (
	"github.com/couchbase/gocb/v2"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("default")

	collection := bucket.DefaultCollection()

	// #tag::rawbinary[]
	input := []byte("hello world")
	transcoder := gocb.NewRawBinaryTranscoder()

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

	var returned []byte
	err = getRes.Content(&returned)
	if err != nil {
		panic(err)
	}
	// #end::rawbinary[]

	if err := cluster.Close(nil); err != nil {
		panic(err)
	}
}
