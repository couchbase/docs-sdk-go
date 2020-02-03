package main

// #tag::connect[]
import (
	"fmt"
	"time"

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
		// handle err
	}
	// #end::connect[]

	// #tag::bucket[]
	// get a bucket reference
	cluster.Bucket("travel-sample")
	// #end::bucket[]

	// #tag::query[]
	results, err := cluster.AnalyticsQuery("SELECT \"hello\" as greeting;", nil)
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

	// always close results and check for errors
	err = results.Close()
	if err != nil {
		panic(err)
	}
	// #end::query[]

	// #tag::simple[]
	results, err = cluster.AnalyticsQuery("select airportname, country from airports where country = 'France';", nil)
	// #end::simple[]

	// #tag::positional[]
	results, err = cluster.AnalyticsQuery(
		"select airportname, country from airports where country = ?;",
		&gocb.AnalyticsOptions{
			PositionalParameters: []interface{}{"France"},
		},
	)
	// #end::positional[]

	// #tag::named[]
	results, err = cluster.AnalyticsQuery(
		"select airportname, country from airports where country = $country;",
		&gocb.AnalyticsOptions{
			NamedParameters: map[string]interface{}{"country": "France"},
		},
	)
	// #end::named[]

	// #tag::options[]
	results, err = cluster.AnalyticsQuery(
		"select airportname, country from airports where country = 'France';",
		&gocb.AnalyticsOptions{
			Priority: true,
			Timeout:  100 * time.Second,
		},
	)
	// #end::options[]

	// #tag::results[]
	results, err = cluster.AnalyticsQuery("select airportname, country from airports where country = 'France';", nil)
	if err != nil {
		panic(err)
	}

	var val interface{}
	for results.Next() {
		err := results.Row(&val)
		if err != nil {
			panic(err)
		}

		fmt.Println(val)
	}

	// always close results and check for errors
	err = results.Close()
	if err != nil {
		panic(err)
	}
	// #end::results[]

	// #tag::metadata[]
	results, err = cluster.AnalyticsQuery("select airportname, country from airports where country = 'France';", nil)
	if err != nil {
		panic(err)
	}

	// we only care about metadata so we can ignore the actual values even though we do need to iterate them first
	for results.Next() {
	}

	err = results.Close()
	if err != nil {
		panic(err)
	}

	metadata, err := results.MetaData()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client context id: %s\n", metadata.ClientContextID)
	fmt.Printf("Elapsed time: %d\n", metadata.Metrics.ElapsedTime)
	fmt.Printf("Execution time: %d\n", metadata.Metrics.ExecutionTime)
	fmt.Printf("Result count: %d\n", metadata.Metrics.ResultCount)
	fmt.Printf("Error count: %d\n", metadata.Metrics.ErrorCount)
	// #end::metadata[]
}
