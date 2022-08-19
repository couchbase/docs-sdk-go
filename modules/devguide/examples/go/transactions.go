package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
}

func initTransactions(cb func(cluster *gocb.Cluster, collection *gocb.Collection)) {
	// tag::init[]
	// Initialize the Couchbase cluster
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

	scope := bucket.Scope("inventory")
	collection := scope.Collection("airport")

	transactions := cluster.Transactions()
	// end::init[]

	throwaway(transactions)

	cb(cluster, collection)
}

func config() {
	// tag::config[]
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
		TransactionsConfig: gocb.TransactionsConfig{
			DurabilityLevel: gocb.DurabilityLevelPersistToMajority,
		},
	}
	// end::config[]

	throwaway(opts)
}

func create() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::create[]
		result, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			// The lambda gets passed an AttemptContext object, which permits getting, inserting,
			// removing and replacing documents, and performing N1QL queries.

			// ... Your transaction logic here ...

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			log.Printf("%+v", err)
		}
		// end::create[]

		throwaway(result)
	})
}

func createSimple() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::create-simple[]
		result, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			if _, err := ctx.Insert(collection, "doc1", map[string]interface{}{}); err != nil {
				return err
			}

			// Replace
			doc, err := ctx.Get(collection, "doc1")
			if err != nil {
				return err
			}

			var content map[string]interface{}
			err = doc.Content(&content)
			if err != nil {
				return err
			}
			content["transactions"] = "are awesome"
			_, err = ctx.Replace(doc, content)
			if err != nil {
				return err
			}

			return nil
		}, nil)
		if err != nil {
			log.Printf("%+v", err)
		}
		// end::create-simple[]

		throwaway(result)
	})
}

func examples() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::examples[]
		scope := cluster.Bucket("travel-sample").Scope("inventory")

		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			// Inserting a doc:
			_, err := ctx.Insert(collection, "doc-a", map[string]interface{}{})
			if err != nil {
				return err
			}

			// Getting documents:
			_, err = ctx.Get(collection, "doc-a")
			// Use err != nil && !errors.Is(err, gocb.ErrDocumentNotFound) if the document may or may not exist
			if err != nil {
				return err
			}

			// Replacing a doc:
			docB, err := ctx.Get(collection, "doc-b")
			if err != nil {
				return err
			}

			var content map[string]interface{}
			err = docB.Content(&content)
			if err != nil {
				return err
			}
			content["transactions"] = "are awesome"
			_, err = ctx.Replace(docB, content)
			if err != nil {
				return err
			}

			// Removing a doc:
			docC, err := ctx.Get(collection, "doc-c")
			if err != nil {
				return err
			}

			err = ctx.Remove(docC)
			if err != nil {
				return err
			}

			// Performing a SELECT N1QL query against a scope:
			qr, err := ctx.Query("SELECT * FROM hotel WHERE country = $1", &gocb.TransactionQueryOptions{
				PositionalParameters: []interface{}{"United Kingdom"},
				Scope:                scope,
			})
			if err != nil {
				return err
			}

			type hotel struct {
				Name string `json:"name"`
			}

			var hotels []hotel
			for qr.Next() {
				var h hotel
				err = qr.Row(&h)
				if err != nil {
					return err
				}

				hotels = append(hotels, h)
			}

			// Performing an UPDATE N1QL query on multiple documents, in the `inventory` scope:
			_, err = ctx.Query("UPDATE route SET airlineid = $1 WHERE airline = $2", &gocb.TransactionQueryOptions{
				PositionalParameters: []interface{}{"airline_137", "AF"},
				Scope:                scope,
			})
			if err != nil {
				return err
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		var ambigErr gocb.TransactionCommitAmbiguousError
		if errors.As(err, &ambigErr) {
			log.Println("Transaction possibly committed")

			log.Printf("%+v", ambigErr)
			return
		}
		var failedErr gocb.TransactionFailedError
		if errors.As(err, &failedErr) {
			log.Println("Transaction did not reach commit point")

			log.Printf("%+v", failedErr)
			return
		}
		if err != nil {
			panic(err)
		}
		// end::examples[]
	})
}

func insert() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::insert[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			_, err := ctx.Insert(collection, "insert-doc", map[string]interface{}{})
			if err != nil {
				return err
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::insert[]
	})
}

func get() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		if _, err := collection.Insert("get-doc", "{}", nil); err != nil {
			panic(err)
		}

		// tag::get[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			doc, err := ctx.Get(collection, "get-doc")
			if err != nil {
				return err
			}

			var content interface{}
			err = doc.Content(&content)
			if err != nil {
				return err
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::get[]

		// tag::getOpt[]
		_, err = cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			doc, err := ctx.Get(collection, "get-doc")
			if err != nil && !errors.Is(err, gocb.ErrDocumentNotFound) {
				return err
			}

			fmt.Println(doc != nil)

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::getOpt[]
	})
}

func getReadOwnWrites() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::getReadOwnWrites[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			_, err := ctx.Insert(collection, "ownwritesdoc", map[string]interface{}{})
			if err != nil {
				return err
			}

			doc, err := ctx.Get(collection, "ownwritesdoc")
			if err != nil {
				return err
			}

			var content interface{}
			err = doc.Content(&content)
			if err != nil {
				return err
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::getReadOwnWrites[]
	})
}

func replace() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		if _, err := collection.Insert("replace-doc", "{}", nil); err != nil {
			panic(err)
		}

		// tag::replace[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			doc, err := ctx.Get(collection, "replace-doc")
			if err != nil {
				return err
			}

			var content map[string]interface{}
			err = doc.Content(&content)
			if err != nil {
				return err
			}
			content["transactions"] = "are awesome"

			_, err = ctx.Replace(doc, content)
			if err != nil {
				return err
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::replace[]
	})
}

func remove() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		if _, err := collection.Insert("remove-doc", "{}", nil); err != nil {
			panic(err)
		}

		// tag::remove[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			doc, err := ctx.Get(collection, "remove-doc")
			if err != nil {
				return err
			}

			err = ctx.Remove(doc)
			if err != nil {
				return err
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::remove[]
	})
}

func querySelect() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::querySelect[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			qr, err := ctx.Query("SELECT * FROM `travel-sample`.inventory.hotel WHERE country = $1", &gocb.TransactionQueryOptions{
				PositionalParameters: []interface{}{"United Kingdom"},
			})
			if err != nil {
				return err
			}

			type hotel struct {
				Name string `json:"name"`
			}

			var hotels []hotel
			for qr.Next() {
				var h hotel
				err = qr.Row(&h)
				if err != nil {
					return err
				}

				hotels = append(hotels, h)
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::querySelect[]
	})
}

func querySelectScope() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::querySelectScope[]
		bucket := cluster.Bucket("travel-sample")
		scope := bucket.Scope("inventory")
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			qr, err := ctx.Query("SELECT * FROM hotel WHERE country = $1", &gocb.TransactionQueryOptions{
				PositionalParameters: []interface{}{"United Kingdom"},
				Scope:                scope,
			})
			if err != nil {
				return err
			}

			type hotel struct {
				Name string `json:"name"`
			}

			var hotels []hotel
			for qr.Next() {
				var h hotel
				err = qr.Row(&h)
				if err != nil {
					return err
				}

				hotels = append(hotels, h)
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::querySelectScope[]
	})
}

func queryUpdate() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::queryUpdate[]
		bucket := cluster.Bucket("travel-sample")
		scope := bucket.Scope("inventory")
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			qr, err := ctx.Query("UPDATE hotel SET price = $1 WHERE url LIKE $2 AND country = $3", &gocb.TransactionQueryOptions{
				PositionalParameters: []interface{}{99.99, "http://marriot%", "United Kingdom"},
				Scope:                scope,
			})
			if err != nil {
				return err
			}

			meta, err := qr.MetaData()
			if err != nil {
				return err
			}

			if meta.Metrics.MutationCount != 1 {
				panic("Should have received 1 mutation")
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::queryUpdate[]
	})
}

func queryComplex() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::queryComplex[]
		bucket := cluster.Bucket("travel-sample")
		scope := bucket.Scope("inventory")
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			// Find all hotels of the chain
			qr, err := ctx.Query("SELECT reviews FROM hotel WHERE url LIKE $1 AND country = $2", &gocb.TransactionQueryOptions{
				PositionalParameters: []interface{}{"http://marriot%", "United Kingdom"},
				Scope:                scope,
			})
			if err != nil {
				return err
			}

			// This function (not provided here) will use a trained machine learning model to provide a
			// suitable price based on recent customer reviews
			updatedPrice := priceFromRecentReviews(qr)

			_, err = ctx.Query("UPDATE hotel SET price = $1 WHERE url LIKE $2 AND country = $3", &gocb.TransactionQueryOptions{
				PositionalParameters: []interface{}{updatedPrice, "http://marriot%", "United Kingdom"},
				Scope:                scope,
			})
			if err != nil {
				return err
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::queryComplex[]
	})
}

func queryInsert() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::queryInsert[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			_, err := ctx.Query("INSERT INTO `default` VALUES ('doc', {'hello':'world'})", nil) // <1>
			if err != nil {
				return err
			}

			st := "SELECT `default`.* FROM `default` WHERE META().id = 'doc'" // <2>
			qr, err := ctx.Query(st, nil)
			if err != nil {
				return err
			}

			meta, err := qr.MetaData()
			if err != nil {
				return err
			}

			if meta.Metrics.ResultCount != 1 {
				panic("Should have received 1 result")
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::queryInsert[]
	})
}

func queryRyow() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::queryRyow[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			_, err := ctx.Insert(collection, "queryRyow", map[string]interface{}{"hello": "world"}) // <1>
			if err != nil {
				return err
			}

			st := "SELECT `default`.* FROM `default` WHERE META().id = 'queryRyow'" // <2>
			qr, err := ctx.Query(st, nil)
			if err != nil {
				return err
			}

			meta, err := qr.MetaData()
			if err != nil {
				return err
			}

			if meta.Metrics.ResultCount != 1 {
				panic("Should have received 1 result")
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::queryRyow[]
	})
}

func queryOptions() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::queryOptions[]
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			_, err := ctx.Query("INSERT INTO `default` VALUES ('queryOpts', {'hello':'world'})",
				&gocb.TransactionQueryOptions{Profile: gocb.QueryProfileModeTimings},
			)
			if err != nil {
				return err
			}

			// There is no commit call, by not returning an error the transaction will automatically commit
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// end::queryOptions[]
	})
}

func calculateLevelForExperience(experience int) int {
	return experience / 10
}

// tag::full[]
func playerHitsMonster(damage int, playerID, monsterID string) {
	type monster struct {
		Hitpoints            int `json:"hitpoints"`
		ExperienceWhenKilled int `json:"experience_when_killed"`
	}

	type player struct {
		Experience int `json:"experience"`
		Level      int `json:"level"`
	}

	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			monsterDoc, err := ctx.Get(collection, monsterID)
			if err != nil {
				return err
			}
			playerDoc, err := ctx.Get(collection, monsterID)
			if err != nil {
				return err
			}

			var monsterContent monster
			if err := monsterDoc.Content(&monsterContent); err != nil {
				return err
			}

			monsterNewHitPoints := monsterContent.Hitpoints - damage

			if monsterNewHitPoints <= 0 {
				// Monster is killed. The remove is just for demoing, and a more realistic
				// example would set a "dead" flag or similar.
				err = ctx.Remove(monsterDoc)
				if err != nil {
					return err
				}

				var playerContent player
				if err := playerDoc.Content(&playerContent); err != nil {
					return err
				}

				// The player earns experience for killing the monster
				playerNewExperience := playerContent.Experience + monsterContent.ExperienceWhenKilled
				playerNewLevel := calculateLevelForExperience(playerNewExperience)

				playerContent.Experience = playerNewExperience
				playerContent.Level = playerNewLevel

				_, err = ctx.Replace(playerDoc, playerContent)
				if err != nil {
					return err
				}
			} else {
				// Monster is damaged but still alive
				monsterContent.Hitpoints = monsterNewHitPoints

				_, err = ctx.Replace(monsterDoc, monsterContent)
				if err != nil {
					return err
				}
			}

			return nil
		}, nil)
		var transactionFailedErr gocb.TransactionFailedError
		if errors.As(err, &transactionFailedErr) {
			// The operation failed. Both the monster and the player will be untouched.

			// Situations that can cause this would include either the monster
			// or player not existing (as get is used), or a persistent
			// failure to be able to commit the transaction, for example on
			// prolonged node failure.
			return
		}

		if err != nil {
			panic(err)
		}
	})
}

// end::full[]

func rollbackCause() {
	type customer struct {
		Balance int `json:"balance"`
	}
	costOfItem := 10

	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::rollbackCause[]
		var ErrBalanceInsufficient = errors.New("insufficient funds")

		_, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			doc, err := ctx.Get(collection, "customer-name")
			if err != nil {
				return err
			}

			var cust customer
			err = doc.Content(&cust)
			if err != nil {
				return err
			}

			if cust.Balance < costOfItem {
				return ErrBalanceInsufficient
			}
			// else continue transaction

			return nil
		}, nil)
		var ambigErr gocb.TransactionCommitAmbiguousError
		if errors.As(err, &ambigErr) {
			// This error can only be thrown at the commit point, after the
			// BalanceInsufficient logic has been passed, so there is no need to
			// check getCause here.
			fmt.Println("Transaction possibly committed")
			fmt.Printf("%+v", ambigErr)
			return
		}

		var transactionFailedErr gocb.TransactionFailedError
		if errors.As(err, &transactionFailedErr) {
			if errors.Is(transactionFailedErr, ErrBalanceInsufficient) {
				// Re-raise the error
				panic(transactionFailedErr)
			} else {
				fmt.Println("Transaction did not reach commit point")
				fmt.Printf("%+v", transactionFailedErr)
			}
			return
		}
		// end::rollbackCause[]
	})
}

func configExpiration() {
	// tag::configExpiration[]
	cluster, err := gocb.Connect("localhost", gocb.ClusterOptions{
		TransactionsConfig: gocb.TransactionsConfig{
			Timeout: 120 * time.Second,
		},
	})
	// end::configExpiration[]
	throwaway(err)
	throwaway(cluster)
}

func customMetadata() {
	// tag::customMetadata[]
	cluster, err := gocb.Connect("localhost", gocb.ClusterOptions{
		TransactionsConfig: gocb.TransactionsConfig{
			MetadataCollection: &gocb.TransactionKeyspace{
				BucketName:     "travel-sample",
				ScopeName:      "transactions",
				CollectionName: "metadata",
			},
		},
	})
	// end::customMetadata[]
	throwaway(err)
	throwaway(cluster)
}

func customMetadataTxn() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::customMetadataTxn[]
		metaCollection := cluster.Bucket("travel-sample").Scope("transactions").Collection("other-metadata")
		result, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			// ... transactional code here ...
			return nil
		}, &gocb.TransactionOptions{
			MetadataCollection: metaCollection,
		})
		// end::customMetadataTxn[]
		throwaway(err)
		throwaway(result)
	})
}

func fullErrorHandling() {
	initTransactions(func(cluster *gocb.Cluster, collection *gocb.Collection) {
		// tag::fullErrorHandling[]
		result, err := cluster.Transactions().Run(func(ctx *gocb.TransactionAttemptContext) error {
			// ... transactional code here ...
			return nil
		}, nil)
		var ambigErr gocb.TransactionCommitAmbiguousError
		if errors.As(err, &ambigErr) {
			fmt.Println("Transaction returned TransactionCommitAmbiguous and may have succeeded")

			// Of course, the application will want to use its own logging rather
			// than fmt.Printf
			fmt.Printf("%+v", ambigErr)
			return
		}
		var transactionFailedErr gocb.TransactionFailedError
		if errors.As(err, &transactionFailedErr) {
			// The transaction definitely did not reach commit point
			fmt.Println("Transaction failed with TransactionFailed")
			fmt.Printf("%+v", transactionFailedErr)
			return
		}
		if err != nil {
			panic(err)
		}

		// The transaction definitely reached the commit point. Unstaging
		// the individual documents may or may not have completed
		if !result.UnstagingComplete {
			// In rare cases, the application may require the commit to have
			// completed.  (Recall that the asynchronous cleanup process is
			// still working to complete the commit.)
			// The next step is application-dependent.
		}
		// end::fullErrorHandling[]
	})

}

func priceFromRecentReviews(qe *gocb.TransactionQueryResult) float32 {
	return 1.0
}

// just used so that we can show creation of resources without the linter complaining.
func throwaway(interface{}) {}
