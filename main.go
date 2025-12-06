package main

import (
	"dmd/initial"
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

	r := routing.New(city, asn, &datacenters)

	mux := http.NewServeMux()
	r.Register(mux)

	http.ListenAndServe(":8080", mux)
}
