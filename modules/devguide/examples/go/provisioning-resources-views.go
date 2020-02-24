package main

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
)

func main() {
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

	bucket := cluster.Bucket("default")
	// #tag::viewmgr[]
	viewMgr := bucket.ViewIndexes()
	// #end::viewmgr[]

	createView(viewMgr)
	getView(viewMgr)
	publishView(viewMgr)
	removeView(viewMgr)
}

func createView(viewMgr *gocb.ViewIndexManager) {
	// #tag::createView[]
	designDoc := gocb.DesignDocument{
		Name: "landmarks",
		Views: map[string]gocb.View{
			"by_country": {
				Map:    "function (doc, meta) { if (doc.type == 'landmark') { emit([doc.country, doc.city], null); } }",
				Reduce: nil,
			},
			"by_activity": {
				Map:    "function (doc, meta) { if (doc.type == 'landmark') { emit(doc.activity, null); } }",
				Reduce: "_count",
			},
		},
	}

	err := viewMgr.UpsertDesignDocument(designDoc, gocb.DesignDocumentNamespaceDevelopment, nil)
	if err != nil {
		panic(err)
	}
	// #end::createView[]
}

func getView(viewMgr *gocb.ViewIndexManager) {
	// #tag::getView[]
	ddoc, err := viewMgr.GetDesignDocument("landmarks", gocb.DesignDocumentNamespaceDevelopment, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(ddoc)
	// #end::getView[]
}

func publishView(viewMgr *gocb.ViewIndexManager) {
	// #tag::publishView[]
	err := viewMgr.PublishDesignDocument("landmarks", nil)
	if err != nil {
		panic(err)
	}
	// #end::publishView[]
}

func removeView(viewMgr *gocb.ViewIndexManager) {
	// #tag::removeView[]
	err := viewMgr.DropDesignDocument("landmarks", gocb.DesignDocumentNamespaceProduction, nil)
	if err != nil {
		panic(err)
	}
	// #end::removeView[]
}
