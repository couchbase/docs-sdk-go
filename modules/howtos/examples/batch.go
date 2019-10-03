package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	// #tag::connect[]
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.Connect("10.112.193.101", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("default", &gocb.BucketOptions{})
	collection := bucket.DefaultCollection()
	// #end::connect[]

	// #tag::loadData[]
	numBatches := 8 // number of batches
	type docType struct {
		Name string
		Data interface{}
	}
	sampleName := "beer-sample"
	sampleZip := fmt.Sprintf("/opt/couchbase/samples/%s.zip", sampleName)
	batches := make(map[int][]gocb.BulkOp)
	var numDocs int

	r, err := zip.OpenReader(sampleZip)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for i, f := range r.File {
		// We only want json files from the docs directory.
		if f.FileInfo().IsDir() || !(strings.HasPrefix(f.Name, sampleName+"/docs/") &&
			strings.HasSuffix(f.Name, ".json")) {
			continue
		}

		fileReader, err := f.Open()
		if err != nil {
			panic(err)
		}
		defer fileReader.Close()

		fileContent, err := ioutil.ReadAll(fileReader)
		if err != nil {
			panic(err)
		}

		var docContent interface{}
		err = json.Unmarshal(fileContent, &docContent)
		if err != nil {
			panic(err)
		}

		_, ok := batches[i%numBatches]
		if !ok {
			batches[i%numBatches] = []gocb.BulkOp{}
		}
		batches[i%numBatches] = append(batches[i%numBatches], &gocb.UpsertOp{
			ID:    f.Name,
			Value: docContent,
		})
		numDocs++
	}
	log.Printf("Loaded %d docs\n", numDocs)
	// #end::loadData[]

	// #tag::send[]
	for _, batch := range batches {
		err := collection.Do(batch, nil)
		if err != nil {
			log.Println(err)
		}

		// Be sure to check each individual operation for errors too.
		for _, op := range batch {
			upsertOp := op.(*gocb.UpsertOp)
			if upsertOp.Err != nil {
				log.Println(err)
			}
		}
	}

	cluster.Close(nil)
	log.Println("Completed")
	// #end::send[]
}
