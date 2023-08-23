package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

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
	collection := bucket.Scope("inventory").Collection("airline")

	if err := bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}

	{
		fmt.Println("Example - [mutate-in]")
		// tag::mutate-in[]
		mops := []gocb.MutateInSpec{
			gocb.UpsertSpec("msrp", 18.00, &gocb.UpsertSpecOptions{}),
		}

		_, err := collection.MutateIn("airline_10", mops, &gocb.MutateInOptions{})
		if err != nil {
			panic(err)
		}
		// end::mutate-in[]
	}

	{
		fmt.Println("Example - [lookup-in]")
		// tag::lookup-in[]
		usersCollection := bucket.Scope("tenant_agent_00").Collection("users")
		ops := []gocb.LookupInSpec{
			gocb.GetSpec("credit_cards[0].type", &gocb.GetSpecOptions{}),
			gocb.GetSpec("credit_cards[0].expiration", &gocb.GetSpecOptions{}),
		}

		_, err := usersCollection.LookupIn("1", ops, &gocb.LookupInOptions{})
		if err != nil {
			panic(err)
		}
		// end::lookup-in[]
	}

	{
		fmt.Println("Example - [counters]")
		// tag::counters[]
		counterDocId := "counter-doc"
		// Increment by 1, creating doc if needed.
		// By using `Initial: 1` we set the starting count(non-negative) to 1 if the document needs to be created.
		// If it already exists, the count will increase by the amount provided in the Delta option(i.e 1).
		collection.Binary().Increment(counterDocId, &gocb.IncrementOptions{Initial: 1, Delta: 1})
		// Decrement by 1
		collection.Binary().Decrement(counterDocId, &gocb.DecrementOptions{Delta: 1})
		// Decrement by 5
		collection.Binary().Decrement(counterDocId, &gocb.DecrementOptions{Delta: 5})
		// end::counters[]
	}

	{
		fmt.Println("Example - [counter-increment]")
		// tag::counter-increment[]
		result, err := collection.Get("counter-doc", &gocb.GetOptions{})
		if err != nil {
			panic(err)
		}

		var value int
		if err := result.Content(&value); err != nil {
			panic(err)
		}

		incrementAmnt := 5
		if shouldIncrementValue(value) {
			collection.Replace(
				"counter-doc",
				value+incrementAmnt,
				&gocb.ReplaceOptions{Cas: result.Cas()},
			)
		}
		// end::counter-increment[]

		fmt.Printf("RESULT: %#v\n", value+incrementAmnt)
	}
}

func shouldIncrementValue(value int) bool {
	fmt.Println("Current value:", value)
	if value == 0 {
		return true
	}
	return false
}
