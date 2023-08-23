package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	// #tag::connect[]
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

	// We wait until the bucket is connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}
	// #end::connect[]

	// #tag::workers[]
	type docType struct {
		Name string
		Data interface{}
	}
	concurrency := 24 // number of goroutines
	workChan := make(chan docType, concurrency)
	shutdownChan := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				select {
				case doc := <-workChan:
					_, err := collection.Upsert(doc.Name, doc.Data, nil)
					if err != nil {
						// We could use errgroup or something to error out nicely here.
						log.Println(err)
					}
				case <-shutdownChan:
					wg.Done()
					return
				}
			}
		}()
	}
	// #end::workers[]

	// #tag::loadData[]
	sampleName := "beer-sample"
	sampleZip := fmt.Sprintf("/opt/couchbase/samples/%s.zip", sampleName)

	r, err := zip.OpenReader(sampleZip)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for _, f := range r.File {
		// We only want json files from the docs directory.
		if f.FileInfo().IsDir() || !(strings.HasPrefix(f.Name, sampleName+"/docs/") &&
			strings.HasSuffix(f.Name, ".json")) {
			continue
		}

		fileReader, err := f.Open()
		if err != nil {
			panic(err)
		}

		fileContent, err := ioutil.ReadAll(fileReader)
		if err != nil {
			fileReader.Close()
			panic(err)
		}
		fileReader.Close()

		var docContent interface{}
		err = json.Unmarshal(fileContent, &docContent)
		if err != nil {
			panic(err)
		}

		workChan <- docType{
			Name: f.Name,
			Data: docContent,
		}
	}
	// #end::loadData[]

	// #tag::wait[]
	// Wait for all of the docs to be uploaded.
	for len(workChan) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	// Signal the goroutines to close, this means that we can wait for them to complete any work that they're doing
	// before we actually finish.
	for i := 0; i < concurrency; i++ {
		shutdownChan <- struct{}{}
	}
	wg.Wait()
	cluster.Close(nil)
	log.Println("Completed")
	// #end::wait[]
}
