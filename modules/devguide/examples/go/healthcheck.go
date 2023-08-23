package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
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

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	// #tag::ping[]
	// We'll ping the KV nodes in our cluster.
	pings, err := bucket.Ping(&gocb.PingOptions{
		ReportID:     "my-report",                                  // <.>
		ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeKeyValue}, // <.>
	})
	if err != nil {
		panic(err)
	}

	for service, pingReports := range pings.Services {
		if service != gocb.ServiceTypeKeyValue {
			panic("we got a service type that we didn't ask for!")
		}

		for _, pingReport := range pingReports {
			if pingReport.State != gocb.PingStateOk {
				fmt.Printf(
					"Node %s at remote %s is not OK, error: %s, latency: %s\n",
					pingReport.ID, pingReport.Remote, pingReport.Error, pingReport.Latency.String(),
				)
			} else {
				fmt.Printf(
					"Node %s at remote %s is OK, latency: %s\n",
					pingReport.ID, pingReport.Remote, pingReport.Latency.String(),
				)
			}
		}
	}

	b, err := json.Marshal(pings) // <.>
	if err != nil {
		panic(err)
	}

	fmt.Printf("Ping report JSON: %s", string(b))
	// #end::ping[]

	// #tag::diagnostics[]
	diagnostics, err := cluster.Diagnostics(&gocb.DiagnosticsOptions{
		ReportID: "my-report", // <.>
	})
	if err != nil {
		panic(err)
	}

	if diagnostics.State != gocb.ClusterStateOnline {
		log.Printf("Overall cluster state is not online\n")
	} else {
		log.Printf("Overall cluster state is online\n")
	}

	for serviceType, diagReports := range diagnostics.Services {
		for _, diagReport := range diagReports {
			if diagReport.State != gocb.EndpointStateConnected {
				fmt.Printf(
					"Node %s at remote %s is not connected on service %s, activity last seen at: %s\n",
					diagReport.ID, diagReport.Remote, serviceType, diagReport.LastActivity.String(),
				)
			} else {
				fmt.Printf(
					"Node %s at remote %s is connected on service %s, activity last seen at: %s\n",
					diagReport.ID, diagReport.Remote, serviceType, diagReport.LastActivity.String(),
				)
			}
		}
	}

	db, err := json.Marshal(diagnostics) // <.>
	if err != nil {
		panic(err)
	}

	fmt.Printf("Diagnostics report JSON: %s", string(db))
	// #end::diagnostics[]
	if err := cluster.Close(nil); err != nil {
		panic(err)
	}

}
