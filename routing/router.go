package routing

import (
	"context"
	"dmd/bots"
	"dmd/initial"
	"dmd/live"
	"dmd/models"
	"dmd/output"
	"dmd/steps"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/oschwald/maxminddb-golang/v2"
	"github.com/redis/go-redis/v9"
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
	now := time.Now()

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

	eventID, clientID, store := ev.Event.ID, ev.Event.ClientID, ev.Init.Data.Shop.MyShopifyDomain

	uaInfo := steps.ParseUA(ev.Event.Context.Navigator.UserAgent)
	ip, ipHash := steps.GetClientIP(req)
	ref := steps.ParseReferrer(ev.Event.Context.Document.Referrer)
	utm, other := steps.ParseUTM(ev.Event.Context.Document.Location.Search), steps.ParseNonUTMParams(ev.Event.Context.Document.Location.Search)
	geo := steps.ExtractGeo(ip, r.City, r.ASN)
	screen := steps.BucketScreenSizes(ev.Event.Context.Window.InnerWidth, ev.Event.Context.Window.InnerHeight, ev.Event.Context.Window.Screen.Width, ev.Event.Context.Window.Screen.Height)
	pageType := steps.Classify(ev.Event.Context.Document.Location.Href)

	geo.DataCenter = bots.FromDataCenter(geo.ASN, r.DataCenters.Orgs)
	genericEval := bots.ExtractSignals(req, ev.Event.Context.Document.Referrer, ev.Event.Context.Navigator.UserAgent)
	specificEval := bots.EvaluateSpecific(req, ev.Event.Context.Document.Referrer, ev.Event.Context.Navigator, ev.Event.Context.Window.InnerWidth, ev.Event.Context.Window.InnerHeight, ev.Event.Context.Window.Screen.Width, ev.Event.Context.Window.Screen.Height, ev.Init.Data.Shop.MyShopifyDomain)
	botScore := bots.EvaluateBot(genericEval, specificEval, geo.DataCenter != "", uaInfo.IsBot)

	sessionStruct := live.CreateSessionStruct(ev, geo, uaInfo, utm, pageType, botScore, ref, param)
	sessionResults, duplicate, err := live.MainLiveWork(r.RDB, sessionStruct, eventID, clientID, store, ev.Event.Context.Document.Location.Href, param, &ev)
	if err != nil {
		http.Error(w, "invalid sessiontmp", http.StatusBadRequest)
		return
	} else if duplicate {
		http.Error(w, "duplicate", http.StatusBadRequest)
		return
	}

	outPut := models.Output{
		EventName:        param,
		EventID:          ev.Event.ID,
		EventTime:        ev.Time,
		TimeIn:           now.Unix(),
		TimeOut:          time.Now().Unix(),
		ShopName:         ev.Init.Data.Shop.Name,
		ShopDomain:       store,
		ClientID:         clientID,
		SessionID:        sessionResults.SessionID,
		CustomerID:       "", // figure out
		LoggedIn:         sessionStruct.IsLoggedIn,
		IP:               ip,
		IPHash:           ipHash,
		SessionStatus:    sessionResults.Status,
		PageType:         string(pageType),
		PreviousPurchase: sessionStruct.HasPreviousPurchase,
		BotScore:         int(botScore),
		Params:           other,
		UA:               uaInfo,
		Geo:              geo,
		UTM:              utm,
		Referrer:         ref,
		Screen:           screen,
		RequestSignals:   genericEval,
		BotSignals:       specificEval,

		RawShopify: ev,
	}

	if err := output.PublishOutput(context.Background(), r.PubSub, "topic", outPut); err != nil {
		http.Error(w, "couldn't push to pub sub", http.StatusBadRequest)
		return
	}

	w.Write([]byte(sessionStruct.City + ip + ipHash + ref.DomainOnly + utm.Campaign + other["a"] + geo.ASNOrg + screen.ScreenHeightBucket + string(pageType) + strconv.Itoa(int(botScore))))
}
