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
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		// handle err
	}
	// #end::connect[]

	// #tag::bucket[]
	// get a bucket reference
	bucket := cluster.Bucket("bucket-name", &gocb.BucketOptions{})
	// #end::bucket[]

	// #tag::collection[]
	// get a collection reference
	collection := bucket.DefaultCollection(&gocb.CollectionOptions{})
	// for a named collection and scope
	// collection := bucket.Scope("my-scope").Collection("my-collection", &gocb.CollectionOptions{})
	// #end::collection[]

	// #tag::typeassert[]
	_, err = collection.Get("key", nil)
	if err != nil {
		if kvError, ok := err.(gocb.KeyValueError); ok {
			fmt.Println(kvError.ID())         // the id of the document
			fmt.Println(kvError.StatusCode()) // the memcached error code
			fmt.Println(kvError.Opaque())     // the unique identifier for the operation
			if kvError.StatusCode() == 0x01 {
				fmt.Println("Document could not be found") // maybe do something like return a 404 to your user
			}
		} else {
			fmt.Printf("An unknown error occurred: %v", err)
		}
	}
	// #end::typeassert[]

	// #tag::helperfunc[]
	_, err = collection.Get("key", nil)
	if err != nil {
		if gocb.IsKeyNotFoundError(err) {
			fmt.Println("Document could not be found") // maybe do something like return a 404 to your user
		} else {
			fmt.Printf("An unknown error occurred: %v", err)
		}
	}
	// #end::helperfunc[]

	// #tag::retries[]
	var changeEmail func(maxRetries int) error
	changeEmail = func(maxRetries int) error {
		result, err := collection.Get("doc-id", nil)
		if err != nil {
			return err
		}

		doc := struct {
			Email string `json:"email"`
		}{}
		err = result.Content(&doc)
		if err != nil {
			return err
		}

		doc.Email = "john.smith@couchbase.com"
		_, err = collection.Replace("doc-id", doc, &gocb.ReplaceOptions{
			Cas: result.Cas(),
		})
		if err == nil {
			return nil
		}

		if gocb.IsRetryableError(err) {
			// IsRetryableError will be true for transient errors, such as a CAS mismatch (indicating
			// another agent concurrently modified the document), or a temporary failure (indicating
			// the server is temporarily unavailable or overloaded).  The operation may or may not
			// have been written, but since it is idempotent we can simply retry it.
			if maxRetries > 0 {
				fmt.Printf("Retrying operation on retryable err %v\n", err)
				return changeEmail(maxRetries)
			}

			// Errors can be transient but still exceed our SLA.
			return fmt.Errorf("maximum retries exceeded, aborting on err %v", err)
		}

		// If the err is not isRetryable, there is perhaps a more permanent or serious error,
		// such as a network failure.
		return fmt.Errorf("aborting on err %v", err)
	}

	err = changeEmail(5)
	if err != nil {
		// What to do here is highly application dependent.  Options could include:
		// - Returning a "please try again later" error back to the end-user (if any)
		// - Logging it for manual human review, and possible follow-up with the end-user (if any)
		fmt.Printf("failed to change email: %v", err)
	}
	// #end::retries[]

	// #tag::queryerror[]
	_, err = cluster.Query("select * from `mybucket`", nil)
	if err != nil {
		if queryErr, ok := err.(gocb.QueryError); ok {
			fmt.Println(queryErr.HTTPStatus()) // the HTTP Status from the server
			fmt.Println(queryErr.ContextID())  // the identifier for the query
			fmt.Println(queryErr.Endpoint())   // the http endpoint used for the query
			fmt.Println(queryErr.Code())       // the error code returned by the server
			fmt.Println(queryErr.Message())    // the error message returned by the server
		}
	}
	// #end::queryerror[]

	// #tag::analyticserror[]
	_, err = cluster.AnalyticsQuery("select * from `mybucket`", nil)
	if err != nil {
		if queryErr, ok := err.(gocb.AnalyticsQueryError); ok {
			fmt.Println(queryErr.HTTPStatus()) // the HTTP Status from the server
			fmt.Println(queryErr.ContextID())  // the identifier for the query
			fmt.Println(queryErr.Endpoint())   // the http endpoint used for the query
			fmt.Println(queryErr.Code())       // the error code returned by the server
			fmt.Println(queryErr.Message())    // the error message returned by the server
		}
	}
	// #end::analyticserror[]

	// #tag::searcherror[]
	maybePrintSearchError := func(err error) {
		if err == nil {
			return
		}
		if searchErrs, ok := err.(gocb.SearchErrors); ok {
			fmt.Println(searchErrs.HTTPStatus()) // the HTTP Status from the server
			fmt.Println(searchErrs.ContextID())  // the identifier for the query
			fmt.Println(searchErrs.Endpoint())   // the http endpoint used for the query
			for _, searchErr := range searchErrs.Errors() {
				fmt.Println(searchErr.Message()) // the error message
			}
		}
	}
	indexName := "my-index"
	query := gocb.SearchQuery{Name: indexName, Query: gocb.NewMatchQuery("matchme")}
	searchResult, err := cluster.SearchQuery(query, nil)
	maybePrintSearchError(err)

	var searchRow gocb.SearchResultHit
	for searchResult.Next(&searchRow) {
		// ...
	}

	err = searchResult.Close()
	maybePrintSearchError(err)
	// #end::searcherror[]

	// #tag::viewerror[]
	maybePrintViewError := func(err error) {
		if err == nil {
			return
		}
		if viewErrs, ok := err.(gocb.ViewQueryErrors); ok {
			fmt.Println(viewErrs.HTTPStatus()) // the HTTP Status from the server
			fmt.Println(viewErrs.Endpoint())   // the http endpoint used for the query
			for _, viewErr := range viewErrs.Errors() {
				fmt.Println(viewErr.Message()) // the error message
				fmt.Println(viewErr.Reason())  // the error reason
			}
		}
	}
	viewResult, err := bucket.ViewQuery("test", "test", nil)
	maybePrintViewError(err)

	var viewRow gocb.ViewRow
	for viewResult.Next(&viewRow) {
		// ...
	}

	err = viewResult.Close()
	maybePrintViewError(err)
	// #end::viewerror[]
}
