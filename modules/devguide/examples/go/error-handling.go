package main

import (
	"errors"
	"fmt"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("travel-sample")
	collection := bucket.Scope("inventory").Collection("airport")

	err = bucket.WaitUntilReady(2*time.Second, nil)
	if err != nil {
		panic(err)
	}

	errorsAs(collection)
	errorsIs(collection)
	errDocumentNotFound(collection)
	errDocumentExists(collection)
	errCasMismatch(collection)
	errDurabilityAmbiguous(collection)
	realWorldErrHandling(collection)
	queryError(cluster)
}

func errorsAs(collection *gocb.Collection) {
	// #tag::as[]
	_, err := collection.Get("key", nil)
	if err != nil {
		var kvError *gocb.KeyValueError
		if errors.As(err, &kvError) {
			fmt.Println(kvError.StatusCode) // the memcached error code
			fmt.Println(kvError.Opaque)     // the unique identifier for the operation
			if kvError.StatusCode == 0x01 {
				fmt.Println("Document could not be found") // maybe do something like return a 404 to your user
			}
		} else {
			fmt.Printf("An unknown error occurred: %v", err)
		}
	}
	// #end::as[]
}

func errorsIs(collection *gocb.Collection) {
	// #tag::is[]
	_, err := collection.Get("does-not-exist", nil)
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			fmt.Println("Document could not be found")
		} else {
			fmt.Printf("An unknown error occurred: %v", err)
		}
	}
	// #end::is[]
}

func errDocumentNotFound(collection *gocb.Collection) {
	// #tag::replace[]
	doc := struct{ Foo string }{Foo: "baz"}
	_, err := collection.Replace("does-not-exist", doc, nil)
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			fmt.Println("Key not be found")
		} else {
			fmt.Printf("An unknown error occurred: %v", err)
		}
	}
	// #end::replace[]
}

func errDocumentExists(collection *gocb.Collection) {
	// #tag::exists[]
	doc := struct{ Foo string }{Foo: "baz"}
	_, err := collection.Insert("does-already-exist", doc, nil)
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentExists) {
			fmt.Println("Key already exists")
		} else {
			fmt.Printf("An unknown error occurred: %v", err)
		}
	}
	// #end::exists[]
}

func errCasMismatch(collection *gocb.Collection) {
	newDoc := struct{}{}
	// #tag::cas[]
	var doOperation func(maxAttempts int) (*gocb.MutationResult, error)
	doOperation = func(maxAttempts int) (*gocb.MutationResult, error) {
		doc, err := collection.Get("doc", nil)
		if err != nil {
			return nil, err
		}

		result, err := collection.Replace("doc", newDoc, &gocb.ReplaceOptions{Cas: doc.Cas()})
		if err != nil {
			if errors.Is(err, gocb.ErrCasMismatch) {
				// Simply recursively retry until maxAttempts is hit
				if maxAttempts == 0 {
					return nil, err
				}
				return doOperation(maxAttempts - 1)
			} else {
				return nil, err
			}
		}

		return result, nil
	}
	// #end::cas[]
}

func errDurabilityAmbiguous(collection *gocb.Collection) {
	// #tag::insert[]
	var doInsert func(docId string, doc []byte, maxAttempts int) (string, error)
	doInsert = func(docId string, doc []byte, maxAttempts int) (string, error) {
		_, err := collection.Insert(docId, doc, &gocb.InsertOptions{
			DurabilityLevel: gocb.DurabilityLevelMajority,
		})
		if err != nil {
			if errors.Is(err, gocb.ErrDocumentExists) {
				// The logic here is that if we failed to insert on the first attempt then
				// it's a true error, otherwise we retried due to an ambiguous error, and
				// it's ok to continue as the operation was actually successful.
				if maxAttempts == 0 {
					return "", err
				}

				return "ok!", nil
			} else if errors.Is(err, gocb.ErrDurabilityAmbiguous) {
				if maxAttempts == 0 {
					return "", err
				}

				return doInsert(docId, doc, maxAttempts-1)
			} else if errors.Is(err, gocb.ErrTimeout) {
				if maxAttempts == 0 {
					return "", err
				}

				return doInsert(docId, doc, maxAttempts-1)
			}

			return "", err
		}

		return "ok!", nil
	}
	// #end::insert[]
}

func realWorldErrHandling(collection *gocb.Collection) {
	// #tag::insert-real[]
	var doInsertReal func(docId string, doc []byte, maxAttempts int, delay time.Duration) (string, error)
	doInsertReal = func(docId string, doc []byte, maxAttempts int, delay time.Duration) (string, error) {
		_, err := collection.Insert(docId, doc, &gocb.InsertOptions{DurabilityLevel: gocb.DurabilityLevelMajority})
		if err != nil {
			if errors.Is(err, gocb.ErrDocumentExists) {
				// The logic here is that if we failed to insert on the first attempt then
				// it's a true error, otherwise we retried due to an ambiguous error, and
				// it's ok to continue as the operation was actually successful.
				if maxAttempts == 0 {
					return "", err
				}

				return "ok!", nil
				// Ambiguous errors.  The operation may or may not have succeeded.  For inserts,
				// the insert can be retried, and a DocumentExistsException indicates it was
				// successful.
			} else if errors.Is(err, gocb.ErrDurabilityAmbiguous) || errors.Is(err, gocb.ErrTimeout) ||
				// Temporary/transient errors that are likely to be resolved
				// on a retry.
				errors.Is(err, gocb.ErrTemporaryFailure) || errors.Is(err, gocb.ErrDurableWriteInProgress) ||
				errors.Is(err, gocb.ErrDurableWriteReCommitInProgress) ||
				// These transient errors won't be returned on an insert, but can be used
				// when writing similar wrappers for other mutation operations.
				errors.Is(err, gocb.ErrCasMismatch) {
				if maxAttempts == 0 {
					return "", err
				}

				time.Sleep(delay)
				return doInsertReal(docId, doc, maxAttempts-1, delay*2)
			}

			return "", err
		}

		return "ok!", nil
	}
	// #end::insert-real[]
}

func queryError(cluster *gocb.Cluster) {
	// #tag::query[]
	_, err := cluster.Query("select * from `someotherbucket`", nil)
	if err != nil {
		var queryErr *gocb.QueryError
		if errors.As(err, &queryErr) {
			fmt.Println(queryErr.ClientContextID) // the identifier for the query
			fmt.Println(queryErr.Endpoint)        // the http endpoint used for the query
			fmt.Println(queryErr.Statement)       // the query statement
			fmt.Println(queryErr.Errors)          // a list of errors codes + messages for why the query failed.
		}
	}
	// #end::query[]
}
