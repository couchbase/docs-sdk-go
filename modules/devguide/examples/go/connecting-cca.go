package main

import (
	"crypto/tls"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// #tag::certconnect[]
	// Load the public/private key pair from file
	cert, err := tls.LoadX509KeyPair("mycert.pem", "mykey.pem")
	if err != nil {
		panic(err)
	}

	opts := gocb.ClusterOptions{
		Authenticator: gocb.CertificateAuthenticator{
			ClientCertificate: &cert,
		},
	}
	// Connect to the cluster using certificates and node key, note: couchbases
	cluster, err := gocb.Connect("couchbases://10.112.193.101", opts)
	if err != nil {
		panic(err)
	}
	// #end::certconnect[]

	err = cluster.Close(nil)
	if err != nil {
		panic(err)
	}
}
