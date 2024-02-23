package main

// #tag::connect[]
import (
	"fmt"
	"github.com/couchbase/gocb/v2/vector"
	"time"

	gocb "github.com/couchbase/gocb/v2"
	"github.com/couchbase/gocb/v2/search"
)

// This example requires an index called `travel-sample-index` to be created
// See modules/test/scripts/init-couchbase/init-buckets.sh(line 47)
func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	// #tag::matchquery[]
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}
	// #end::connect[]

	// For Server versions 6.5 or later you do not need to open a bucket here
	bucket := cluster.Bucket("travel-sample")

	// We wait until the bucket is definitely connected and setup.
	// For Server versions 6.5 or later if we hadn't opened a bucket then we could use cluster.WaitUntilReady here.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	matchResult, err := cluster.SearchQuery(
		"travel-sample-index",
		search.NewMatchQuery("swanky"),
		&gocb.SearchOptions{
			Limit:  10,
			Fields: []string{"description"},
		},
	)
	if err != nil {
		panic(err)
	}
	// #end::matchquery[]

	// #tag::iteratingrows[]
	for matchResult.Next() {
		row := matchResult.Row()
		docID := row.ID
		score := row.Score

		var fields interface{}
		err := row.Fields(&fields)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Document ID: %s, search score: %f, fields included in result: %v\n", docID, score, fields)
	}

	// always check for errors after iterating
	err = matchResult.Err()
	if err != nil {
		panic(err)
	}
	// #end::iteratingrows[]

	// #tag::daterangequery[]
	dateRangeResult, err := cluster.SearchQuery(
		"travel-sample-index",
		search.NewDateRangeQuery().Start("2019-01-01", true).End("2019-02-01", false),
		&gocb.SearchOptions{
			Limit: 10,
		},
	)
	if err != nil {
		panic(err)
	}
	// #end::daterangequery[]

	for dateRangeResult.Next() {
		row := dateRangeResult.Row()
		docID := row.ID
		score := row.Score

		var fields interface{}
		if err := row.Fields(&fields); err != nil {
			panic(err)
		}

		fmt.Printf("Document ID: %s, search score: %f, fields included in range result: %v\n", docID, score, fields)
	}

	// always check for errors after iterating
	err = dateRangeResult.Err()
	if err != nil {
		panic(err)
	}

	// #tag::conjunctionquery[]
	conjunctionResult, err := cluster.SearchQuery(
		"travel-sample-index",
		search.NewConjunctionQuery(
			search.NewMatchQuery("swanky"),
			search.NewDateRangeQuery().Start("2019-01-01", true).End("2019-02-01", false),
		),
		&gocb.SearchOptions{
			Limit: 10,
		},
	)
	if err != nil {
		panic(err)
	}
	// #end::conjunctionquery[]

	for conjunctionResult.Next() {
		row := conjunctionResult.Row()
		docID := row.ID
		score := row.Score

		var fields interface{}
		err := row.Fields(&fields)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Document ID: %s, search score: %f, fields included in conjunction result: %v\n", docID, score, fields)
	}

	// always check for errors after iterating
	err = conjunctionResult.Err()
	if err != nil {
		panic(err)
	}

	// #tag::iteratingfacets[]
	facetsResult, err := cluster.SearchQuery(
		"travel-sample-index",
		search.NewMatchAllQuery(),
		&gocb.SearchOptions{
			Limit: 10,
			Facets: map[string]search.Facet{
				"type": search.NewTermFacet("type", 5),
			},
		},
	)
	if err != nil {
		panic(err)
	}

	for facetsResult.Next() {
	}

	facets, err := facetsResult.Facets()
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
	collection := bucket.Scope("inventory").Collection("hotel")

	hotel := struct {
		Description string `json:"description"`
	}{Description: "super swanky"}
	myWriteResult, err := collection.Upsert("a-new-hotel", hotel, nil)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)

	consistentWithResult, err := cluster.SearchQuery(
		"travel-sample-index",
		search.NewMatchQuery("swanky"),
		&gocb.SearchOptions{
			Limit:          10,
			ConsistentWith: gocb.NewMutationState(*myWriteResult.MutationToken()),
		},
	)
	if err != nil {
		panic(err)
	}
	// #end::consistency[]

	for consistentWithResult.Next() {
		row := consistentWithResult.Row()
		docID := row.ID
		score := row.Score

		fmt.Printf("Document ID: %s, search score: %f\n", docID, score)
	}

	// always check for errors after iterating
	err = consistentWithResult.Err()
	if err != nil {
		panic(err)
	}

	vectorQuery := []float32{1.1, 1.2}
	anotherVectorQuery := []float32{0.9, 0.2}
	scope := bucket.DefaultScope()

	{
		// #tag::single_vector_query[]
		request := gocb.SearchRequest{
			VectorSearch: vector.NewSearch(
				[]*vector.Query{
					vector.NewQuery("vector_field", vectorQuery),
				}, nil,
			),
		}
		vectorResult, err := scope.Search("vector-index", request, nil)
		if err != nil {
			panic(err)
		}
		// #end::single_vector_query[]

		for vectorResult.Next() {
			row := vectorResult.Row()
			docID := row.ID
			score := row.Score

			fmt.Printf("Document ID: %s, search score: %f\n", docID, score)
		}

		// always check for errors after iterating
		err = vectorResult.Err()
		if err != nil {
			panic(err)
		}
	}
	{
		// #tag::multiple_vector_queries[]
		request := gocb.SearchRequest{
			VectorSearch: vector.NewSearch(
				[]*vector.Query{
					vector.NewQuery("vector_field", vectorQuery).NumCandidates(2).Boost(0.3),
					vector.NewQuery("vector_field", anotherVectorQuery).NumCandidates(5).Boost(0.7),
				},
				&vector.SearchOptions{
					VectorQueryCombination: vector.VectorQueryCombinationAnd,
				},
			),
		}
		vectorResult, err := scope.Search("vector-index", request, nil)
		if err != nil {
			panic(err)
		}
		// #end::multiple_vector_queries[]

		for vectorResult.Next() {
			row := vectorResult.Row()
			docID := row.ID
			score := row.Score

			fmt.Printf("Document ID: %s, search score: %f\n", docID, score)
		}

		// always check for errors after iterating
		err = vectorResult.Err()
		if err != nil {
			panic(err)
		}
	}
	{
		// #tag::vector_fts_query_combination[]
		request := gocb.SearchRequest{
			VectorSearch: vector.NewSearch(
				[]*vector.Query{
					vector.NewQuery("vector_field", vectorQuery).NumCandidates(2).Boost(0.3),
					vector.NewQuery("vector_field", anotherVectorQuery).NumCandidates(5).Boost(0.7),
				}, nil,
			),
			SearchQuery: search.NewMatchAllQuery(),
		}
		vectorResult, err := scope.Search("vector-and-fts-index", request, nil)
		if err != nil {
			panic(err)
		}
		// #end::vector_fts_query_combination[]

		for vectorResult.Next() {
			row := vectorResult.Row()
			docID := row.ID
			score := row.Score

			fmt.Printf("Document ID: %s, search score: %f\n", docID, score)
		}

		// always check for errors after iterating
		err = vectorResult.Err()
		if err != nil {
			panic(err)
		}
	}
	{
		// #tag::fts_search_request[]
		request := gocb.SearchRequest{
			SearchQuery: search.NewMatchQuery("swanky"),
		}

		result, err := scope.Search("travel-sample-index", request, nil)
		if err != nil {
			panic(err)
		}
		// #end::fts_search_request[]

		for result.Next() {
			row := result.Row()
			docID := row.ID
			score := row.Score

			fmt.Printf("Document ID: %s, search score: %f\n", docID, score)
		}

		// always check for errors after iterating
		err = result.Err()
		if err != nil {
			panic(err)
		}
	}
}
