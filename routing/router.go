package routing

import (
	"dmd/initial"
	"dmd/models"
	"dmd/steps"
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

	if !steps.CheckEvent(param) {
		http.Error(w, "invalid event", http.StatusBadRequest)
		return
	}

	var ev models.IngestEvent
	err := json.NewDecoder(req.Body).Decode(&ev)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if !steps.CheckSessionID(ev.Session.Current) {
		http.Error(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	uaInfo := steps.ParseUA(ev.Event.Context.Navigator.UserAgent)
	ip, ipHash := steps.GetClientIP(req)
	ref := steps.ParseReferrer(ev.Event.Context.Document.Referrer)
	utm, other := steps.ParseUTM(ev.Event.Context.Document.Location.Search), steps.ParseNonUTMParams(ev.Event.Context.Document.Location.Search)
	geo := steps.ExtractGeo(ip, r.City, r.ASN)
	screen := steps.BucketScreenSizes(ev.Event.Context.Window.InnerWidth, ev.Event.Context.Window.InnerHeight, ev.Event.Context.Window.Screen.Width, ev.Event.Context.Window.Screen.Height)
	pageType := steps.Classify(ev.Event.Context.Document.Location.Href)

	w.Write([]byte(uaInfo.BotCategory + ip + ipHash + ref.DomainOnly + utm.Campaign + other["a"] + geo.ASNOrg + screen.ScreenHeightBucket + string(pageType)))
}
