package main

import (
	"net/http"
	"os"

	"github.com/oschwald/maxminddb-golang/v2"
)

type DataCenterASNs struct {
	set map[uint]struct{}
}

var (
	asnList DataCenterASNs
	geoDB   *maxminddb.Reader
)

func LoadASNs(path string) (DataCenterASNs, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return DataCenterASNs{}, err
	}
	s := DataCenterASNs{set: map[uint]struct{}{}}
	start := 0
	for i := 0; i <= len(b); i++ {
		if i == len(b) || b[i] == '\n' {
			if i > start {
				var n uint
				parseNum(b[start:i], &n)
				s.set[n] = struct{}{}
			}
			start = i + 1
		}
	}
	return s, nil
}

func parseNum(b []byte, out *uint) {
	var n uint
	for _, c := range b {
		if c >= '0' && c <= '9' {
			n = n*10 + uint(c-'0')
		}
	}
	*out = n
}

func main() {
	a, err := LoadASNs("files/asn.txt")
	if err != nil {
		panic(err)
	}
	asnList = a

	db, err := maxminddb.Open("files/GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}
	geoDB = db

	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.ListenAndServe(":8080", nil)
}
