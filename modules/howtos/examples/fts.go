package main

// #tag::connect[]
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
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}
	// #end::connect[]

	collection := cluster.Bucket("my-bucket", nil).DefaultCollection()

	// #tag::matchquery[]
	matchResult, err := cluster.SearchQuery(
		"travel-sample-index-hotel-description",
		gocb.NewMatchQuery("swanky"),
		&gocb.SearchOptions{
			Limit: 10,
		},
	)
	// #end::matchResult[]
	fmt.Println(matchResult)

	// #tag::daterangequery[]
	dateRangeResult, err := cluster.SearchQuery(
		"travel-sample-index-hotel-description",
		gocb.NewDateRangeQuery().Start("2019-01-01", true).End("2019-02-01", false),
		&gocb.SearchOptions{
			Limit: 10,
		},
	)
	// #end::daterangequery[]
	fmt.Println(dateRangeResult)

	// #tag::conjunctionquery[]
	conjunctionResult, err := cluster.SearchQuery(
		"travel-sample-index-hotel-description",
		gocb.NewConjunctionQuery(
			gocb.NewMatchQuery("swanky"),
			gocb.NewDateRangeQuery().Start("2019-01-01", true).End("2019-02-01", false),
		),
		&gocb.SearchOptions{
			Limit: 10,
		},
	)
	// #end::conjunctionquery[]
	fmt.Println(conjunctionResult)

	// #tag::iteratingrows[]
	var row gocb.SearchRow
	for matchResult.Next(&row) {
		docID := row.ID
		score := row.Score

		var fields interface{}
		err := row.Fields(&fields)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Document ID: %s, search score: %d, fields included in result: %v\n", docID, score, fields)
	}

	err = matchResult.Close()
	if err != nil {
		panic(err)
	}
	// #end::iteratingrows[]

	// #tag::iteratingfacets[]
	facets, err := matchResult.Facets()
	if err != nil {
		panic(err)
	}
	for _, facet := range facets {
		field := facet.Field
		total := facet.Total

		fmt.Printf("Facet field: %s, total: %d\n", field, total)
	}
	// #end::iteratingfacets[]

	// #tag::consistency[]
	hotel := struct {
		Description string `json:"description"`
	}{Description: "super swanky"}
	myWriteResult, err := collection.Upsert("a-new-hotel", hotel, nil)
	if err != nil {
		panic(err)
	}

	consistentWithResult, err := cluster.SearchQuery(
		"travel-sample-index-hotel-description",
		gocb.NewMatchQuery("swanky"),
		&gocb.SearchOptions{
			Limit:          10,
			ConsistentWith: gocb.NewMutationState(*myWriteResult.MutationToken()),
		},
	)
	// #end::consistency[]
	fmt.Println(consistentWithResult)
}
