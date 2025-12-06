package routing

import (
	"dmd/initial"
	"net/http"

	"github.com/oschwald/maxminddb-golang/v2"
)

type Router struct {
	City        *maxminddb.Reader
	ASN         *maxminddb.Reader
	DataCenters *initial.DataCenterASNs
}

func New(a *maxminddb.Reader, b *maxminddb.Reader, c *initial.DataCenterASNs) *Router {
	return &Router{City: a, ASN: b, DataCenters: c}
}

func (r *Router) Register(mux *http.ServeMux) {
	mux.HandleFunc("/ingest/", func(w http.ResponseWriter, req *http.Request) {
		param := req.URL.Path[len("/ingest/"):]
		r.Ingest(w, req, param)
	})
}

func (r *Router) Ingest(w http.ResponseWriter, req *http.Request, param string) {
	w.WriteHeader(200)
}
