package main

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// tag::change-password[]
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
	collection := bucket.Scope("inventory").Collection("airline")

	// Change the current user's password.
	userMgr := cluster.Users()

	newPassword := "newpassword"
	if err := userMgr.ChangePassword(newPassword, &gocb.ChangePasswordOptions{}); err != nil {
		panic(err)
	}

	// Reconnect your client
	opts = gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			// Use the new password
			Password: newPassword,
		},
	}
	cluster, err = gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully changed the user's password")

	// Perform an operation with the newly authenticated user
	_, err = collection.Get("airline_10", nil)
	if err != nil {
		panic(err)
	}
	// end::change-password[]

	resetPasswordToDefault(cluster)
}

// Resets the user's password to its default value, to clear up example changes.
func resetPasswordToDefault(cluster *gocb.Cluster) {
	userMgr := cluster.Users()

	if err := userMgr.ChangePassword("password", &gocb.ChangePasswordOptions{}); err != nil {
		panic(err)
	}

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

	fmt.Println("Successfully reset password to default value")
}
