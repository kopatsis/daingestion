package routing

import (
	"context"
	"dmd/bots"
	"dmd/initial"
	"dmd/models"
	"dmd/output"
	"dmd/steps"
	"encoding/json"
	"net/http"
	"strconv"

	"cloud.google.com/go/pubsub/v2"
	"github.com/go-redis/redis"
	"github.com/oschwald/maxminddb-golang/v2"
)

type Router struct {
	City        *maxminddb.Reader
	ASN         *maxminddb.Reader
	DataCenters *initial.DataCenterASNs
	PubSub      *pubsub.Client
	RDB         *redis.Client
}

func New(a *maxminddb.Reader, b *maxminddb.Reader, c *initial.DataCenterASNs, d *pubsub.Client, e *redis.Client) *Router {
	return &Router{City: a, ASN: b, DataCenters: c, PubSub: d, RDB: e}
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

	datacenter := bots.FromDataCenter(geo.ASN, r.DataCenters.Orgs)
	genericEval := bots.ExtractSignals(req, ev.Event.Context.Document.Referrer, ev.Event.Context.Navigator.UserAgent)
	specificEval := bots.EvaluateSpecific(req, ev.Event.Context.Document.Referrer, ev.Event.Context.Navigator, ev.Event.Context.Window.InnerWidth, ev.Event.Context.Window.InnerHeight, ev.Event.Context.Window.Screen.Width, ev.Event.Context.Window.Screen.Height, ev.Init.Data.Shop.MyShopifyDomain)
	botScore := bots.EvaluateBot(genericEval, specificEval, datacenter != "", uaInfo.IsBot)

	outputData := output.Output{EventName: param, Timestamp: ev.Event.Timestamp, Data: ""}

	if err := output.PublishOutput(context.Background(), r.PubSub, "topic", outputData); err != nil {
		http.Error(w, "couldn't push to pub sub", http.StatusBadRequest)
		return
	}

	w.Write([]byte(uaInfo.BotCategory + ip + ipHash + ref.DomainOnly + utm.Campaign + other["a"] + geo.ASNOrg + screen.ScreenHeightBucket + string(pageType) + datacenter + strconv.Itoa(int(botScore))))
}
