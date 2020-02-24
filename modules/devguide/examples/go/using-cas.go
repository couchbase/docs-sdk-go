package main

import (
	"errors"
	"time"

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
	collection := bucket.DefaultCollection()

	replaceWithCas(collection, "userID")
	lockingAndCas(collection)

	cluster.Close(nil)
}

// #tag::handlingerrors[]
type user struct {
	visitCount int
}

func replaceWithCas(collection *gocb.Collection, userID string) {
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		// Get the current document contents
		getRes, err := collection.Get(userID, nil)
		if err != nil {
			panic(err)
		}

		// Increment the visit count
		var userDoc user
		err = getRes.Content(&userDoc)
		if err != nil {
			panic(err)
		}
		userDoc.visitCount++

		// Attempt to replace the document using CAS
		_, err = collection.Replace(userID, userDoc, &gocb.ReplaceOptions{
			Cas: getRes.Cas(),
		})
		if err != nil {
			// Check if the error thrown is a cas mismatch, if it is, we retry
			if errors.Is(err, gocb.ErrCasMismatch) {
				continue
			}
			panic(err)
		}

		// If no errors occured during the replace, we can exit our retry loop
		break
	}
}

// #end::handlingerrors[]

func lockingAndCas(collection *gocb.Collection) {
	// #tag::locking[]
	getRes, err := collection.GetAndLock("key", 2*time.Second, nil)
	if err != nil {
		panic(err)
	}

	lockedCas := getRes.Cas()

	/* an example of simply unlocking the document:
	collection.Unlock("key", lockedCas, nil)
	*/

	_, err = collection.Replace("key", "new value", &gocb.ReplaceOptions{
		Cas: lockedCas,
	})
	// #end::locking[]
}
