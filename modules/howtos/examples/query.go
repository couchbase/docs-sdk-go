package main

// #tag::connect[]
import (
	"fmt"

	gocb "github.com/couchbase/gocb/v2"
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
	// #end::connect[]

	bucket := cluster.Bucket("travel-sample")
	collection := bucket.DefaultCollection()

	simple(cluster)

	// #tag::pos-params[]
	query := "SELECT x.* FROM `travel-sample` x WHERE x.`type`=$1 LIMIT 10;"
	rows, err := cluster.Query(query, &gocb.QueryOptions{PositionalParameters: []interface{}{"hotel"}})
	// #end::pos-params[]
	fmt.Println(rows)

	// #tag::named-params[]
	query = "SELECT x.* FROM `travel-sample` x WHERE x.`type`=$type LIMIT 10;"
	params := make(map[string]interface{}, 1)
	params["type"] = "hotel"
	rows, err = cluster.Query(query, &gocb.QueryOptions{NamedParameters: params})
	// #end::named-params[]
	fmt.Println(rows)

	// #tag::results[]
	query = "SELECT x.* FROM `travel-sample` x WHERE x.`type`=$1 LIMIT 10;"
	results, err := cluster.Query(query, &gocb.QueryOptions{PositionalParameters: []interface{}{"hotel"}})

	// check query was successful
	if err != nil {
		panic(err)
	}

	// iterate over rows
	for rows.Next() {
		var hotel interface{} // this could also be a specific type like Hotel
		err := rows.Row(&hotel)
		if err != nil {
			panic(err)
		}
		fmt.Println(hotel)
	}

	// always check for errors after iterating
	err = results.Err()
	if err != nil {
		panic(err)
	}
	// #end::results[]

	metrics(results)
	one(cluster)

	// #tag::consistentwith[]
	// create / update document (mutation)
	result, err := collection.Upsert("id", struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}{Name: "somehotel", Type: "hotel"}, nil)
	if err != nil {
		panic(err)
	}

	// create mutation state from mutation results
	state := gocb.NewMutationState(*result.MutationToken())

	// use mutation state with query option
	rows, err = cluster.Query("SELECT x.* FROM `travel-sample` x WHERE x.`type`=\"hotel\" LIMIT 10", &gocb.QueryOptions{
		ConsistentWith: state,
	})
	// #end::consistentwith[]
	fmt.Println(rows)
}

func simple(cluster *gocb.Cluster) {
	// #tag::simple[]
	results, err := cluster.Query("SELECT \"hello\" as greeting;", nil)
	if err != nil {
		panic(err)
	}

	var greeting interface{}
	for results.Next() {
		err := results.Row(&greeting)
		if err != nil {
			panic(err)
		}
		fmt.Println(greeting)
	}

	// always check for errors after iterating
	err = results.Err()
	if err != nil {
		panic(err)
	}
	// #end::simple[]
}

func one(cluster *gocb.Cluster) {
	// #tag::one[]
	query := "SELECT x.* FROM `travel-sample` x WHERE x.`type`=$1 LIMIT 1;"
	results, err := cluster.Query(query, &gocb.QueryOptions{PositionalParameters: []interface{}{"hotel"}})

	// check query was successful
	if err != nil {
		panic(err)
	}

	var hotel interface{} // this could also be a specific type like Hotel
	err = results.One(&hotel)
	if err != nil {
		panic(err)
	}
	fmt.Println(hotel)
	// #end::one[]
}

func metrics(results *gocb.QueryResult) {
	// #tag::metrics[]
	metadata, err := results.MetaData()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Execution Time: %d\n", metadata.Metrics.ExecutionTime)
	// #end::metrics[]
}

func consistency(cluster *gocb.Cluster) {
	// #tag::consistency[]
	query := "SELECT x.* FROM `travel-sample` x WHERE x.`type`=\"hotel\" LIMIT 10"
	results, err := cluster.Query(query, &gocb.QueryOptions{
		ScanConsistency: gocb.QueryScanConsistencyRequestPlus,
	})
	if err != nil {
		panic(err)
	}
	// #end::consistency[]
	fmt.Println(results)
}
