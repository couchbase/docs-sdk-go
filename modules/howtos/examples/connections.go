package main

import (
	gocb "github.com/couchbase/gocb/v2"
)

func simpleconnect() {
	// #tag::simpleconnect[]
	opts := gocb.ClusterOptions{
		Username: "Administrator",
		Password: "password",
	}
	cluster, err := gocb.Connect("10.112.193.101", opts)
	if err != nil {
		panic(err)
	}
	// #end::simpleconnect[]

	cluster.Close()
}

func multinodeconnect() {
	// #tag::multinodeconnect[]
	opts := gocb.ClusterOptions{
		Username: "Administrator",
		Password: "password",
	}
	cluster, err := gocb.Connect("192.168.56.101,192.168.56.102", opts)
	if err != nil {
		panic(err)
	}
	// #end::multinodeconnect[]

	cluster.Close()
}

func customports() {
	// #tag::customports[]
	opts := gocb.ClusterOptions{
		Username: "Administrator",
		Password: "password",
	}
	cluster, err := gocb.Connect("couchbase://192.168.56.101:1234,192.168.56.102:5678", opts)
	if err != nil {
		panic(err)
	}
	// #end::customports[]

	cluster.Close()
}

func tlsconnect() {
	// #tag::tls[]
	// We use the system certificate pool here and assume the Couchbase root certificate(s) have
	// been installed there, but it is also possible to load these from a file.
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}

	opts := gocb.ClusterOptions{
		Username: "Administrator",
		Password: "password",
		SecurityConfig: SecurityConfig{
			TLSRootCAs: rootCAs,
		}
	}
	cluster, err := gocb.Connect("couchbases://10.112.193.101", opts)
	if err != nil {
		panic(err)
	}
	// #end::tls[]

	cluster.Close()
}

func dnssrv() {
	// #tag::dnssrv[]
	opts := gocb.ClusterOptions{
		Username: "Administrator",
		Password: "password",
	}
	cluster, err := gocb.Connect("couchbase://mysrvrecord.hostname.com", opts)
	if err != nil {
		panic(err)
	}
	// #end::dnssrv[]

	cluster.Close()
}

func main() {
	simpleconnect()
	multinodeconnect()
	customports()
	tlsconnect()
	dnssrv()
}
