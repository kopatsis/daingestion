package routing

import (
	"dmd/initial"
	"dmd/models"
	"encoding/json"
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
	valid := map[string]struct{}{
		"cart_viewed":                      {},
		"checkout_address_info_submitted":  {},
		"checkout_completed":               {},
		"checkout_contact_info_submitted":  {},
		"checkout_shipping_info_submitted": {},
		"checkout_started":                 {},
		"collection_viewed":                {},
		"page_viewed":                      {},
		"payment_info_submitted":           {},
		"product_added_to_cart":            {},
		"product_removed_from_cart":        {},
		"product_viewed":                   {},
		"search_submitted":                 {},
	}

	if _, ok := valid[param]; !ok {
		http.Error(w, "invalid event", http.StatusBadRequest)
		return
	}

	var ev models.IngestEvent
	err := json.NewDecoder(req.Body).Decode(&ev)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
}
