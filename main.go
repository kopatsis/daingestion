package main

import (
	"context"
	"dmd/initial"
	"dmd/output"
	"dmd/routing"
	"net/http"
)

func main() {
	datacenters, err := initial.LoadASNs()
	if err != nil {
		panic(err)
	}

	city, asn, err := initial.LoadGeoLite()
	if err != nil {
		panic(err)
	}

	pubsubClient, err := output.NewPubSubClient(context.Background(), "id")
	if err != nil {
		panic(err)
	}

	r := routing.New(city, asn, &datacenters, pubsubClient)

	mux := http.NewServeMux()
	r.Register(mux)

	http.ListenAndServe(":8080", mux)
}
