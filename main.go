package main

import (
	"context"
	"dmd/initial"
	"dmd/logging"
	"dmd/output"
	"dmd/routing"
	"net/http"

	"github.com/gamebtc/devicedetector"
)

func main() {
	datacenters, err := initial.LoadASNs()
	if err != nil {
		logging.LogError(
			"ERROR",
			"asn_datacenter_load_failed",
			"startup",
			"",
			"",
			"",
			false,
			"failed to load asn or datacenter data",
		)
		panic(err)
	}

	city, asn, err := initial.LoadGeoLite()
	if err != nil {
		logging.LogError(
			"ERROR",
			"geolite_load_failed",
			"startup",
			"",
			"",
			"",
			false,
			"failed to load asn or datacenter data",
		)

		panic(err)
	}

	pubsubClient, err := output.NewPubSubClient(context.Background(), "id")
	if err != nil {
		logging.LogError(
			"CRITICAL",
			"redis_startup_failed",
			"startup",
			"",
			"",
			"",
			true,
			"failed to initialize redis client",
		)

		panic(err)
	}

	dd, err := devicedetector.NewDeviceDetector("regexes")
	if err != nil {
		logging.LogError(
			"ERROR",
			"device_detector_load_failed",
			"startup",
			"",
			"",
			"",
			false,
			"failed to load asn or datacenter data",
		)
		panic(err)
	}

	rdb := initial.NewRedis()

	r := routing.New(city, asn, &datacenters, pubsubClient, rdb, dd)

	mux := http.NewServeMux()
	r.Register(mux)

	http.ListenAndServe(":8080", mux)
}
