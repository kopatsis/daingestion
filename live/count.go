package live

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionActiveState struct {
	Country   string
	Region    string
	City      string
	Latitude  float64
	Longitude float64
	Timezone  string

	IsBot               bool
	ASN                 uint
	ASNProvider         string
	IsLoggedIn          bool
	HasPreviousPurchase bool

	DeviceType  string
	DeviceBrand string
	OSName      string
	BrowserName string
	Language    string

	RefDomain string
	UTMSource string

	RouteType string
	Route     string
	FullURL   string

	IsViewingProduct bool
	ActiveProductID  string
	ActiveVariantID  string

	IsViewingCollection bool
	ActiveCollectionID  string

	HasActiveCart bool
	IsViewingCart bool
	IsInCheckout  bool
}

func SetState(ctx context.Context, rdb *redis.Client, sessionID string, s SessionActiveState) error {
	key := "session:" + sessionID
	_, err := rdb.HSet(ctx, key, map[string]interface{}{
		"Country":             s.Country,
		"Region":              s.Region,
		"City":                s.City,
		"Latitude":            s.Latitude,
		"Longitude":           s.Longitude,
		"Timezone":            s.Timezone,
		"IsBot":               strconv.FormatBool(s.IsBot),
		"ASN":                 strconv.Itoa(int(s.ASN)),
		"ASNProvider":         s.ASNProvider,
		"IsLoggedIn":          strconv.FormatBool(s.IsLoggedIn),
		"HasPreviousPurchase": strconv.FormatBool(s.HasPreviousPurchase),
		"DeviceType":          s.DeviceType,
		"DeviceBrand":         s.DeviceBrand,
		"OSName":              s.OSName,
		"BrowserName":         s.BrowserName,
		"Language":            s.Language,
		"RefDomain":           s.RefDomain,
		"UTMSource":           s.UTMSource,
		"RouteType":           s.RouteType,
		"Route":               s.Route,
		"FullURL":             s.FullURL,
		"IsViewingProduct":    strconv.FormatBool(s.IsViewingProduct),
		"ActiveProductID":     s.ActiveProductID,
		"ActiveVariantID":     s.ActiveVariantID,
		"IsViewingCollection": strconv.FormatBool(s.IsViewingCollection),
		"ActiveCollectionID":  s.ActiveCollectionID,
		"HasActiveCart":       strconv.FormatBool(s.HasActiveCart),
		"IsViewingCart":       strconv.FormatBool(s.IsViewingCart),
		"IsInCheckout":        strconv.FormatBool(s.IsInCheckout),
	}).Result()
	if err != nil {
		return err
	}
	return rdb.Expire(ctx, key, 300000000000).Err()
}

func GetState(ctx context.Context, rdb *redis.Client, sessionID string) SessionActiveState {
	key := "session:" + sessionID
	m, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return SessionActiveState{}
	}

	b1, _ := strconv.ParseBool(m["IsBot"])
	asn, _ := strconv.ParseUint(m["ASN"], 10, 32)
	lat, _ := strconv.ParseFloat(m["Latitude"], 64)
	long, _ := strconv.ParseFloat(m["Longitude"], 64)
	b2, _ := strconv.ParseBool(m["IsLoggedIn"])
	b3, _ := strconv.ParseBool(m["HasPreviousPurchase"])
	b4, _ := strconv.ParseBool(m["IsViewingProduct"])
	b5, _ := strconv.ParseBool(m["IsViewingCollection"])
	b6, _ := strconv.ParseBool(m["HasActiveCart"])
	b7, _ := strconv.ParseBool(m["IsViewingCart"])
	b8, _ := strconv.ParseBool(m["IsInCheckout"])

	return SessionActiveState{
		Country:             m["Country"],
		Region:              m["Region"],
		City:                m["City"],
		Latitude:            lat,
		Longitude:           long,
		Timezone:            m["Timezone"],
		IsBot:               b1,
		ASN:                 uint(asn),
		ASNProvider:         m["ASNProvider"],
		IsLoggedIn:          b2,
		HasPreviousPurchase: b3,
		DeviceType:          m["DeviceType"],
		DeviceBrand:         m["DeviceBrand"],
		OSName:              m["OSName"],
		BrowserName:         m["BrowserName"],
		Language:            m["Language"],
		RefDomain:           m["RefDomain"],
		UTMSource:           m["UTMSource"],
		RouteType:           m["RouteType"],
		Route:               m["Route"],
		FullURL:             m["FullURL"],
		IsViewingProduct:    b4,
		ActiveProductID:     m["ActiveProductID"],
		ActiveVariantID:     m["ActiveVariantID"],
		IsViewingCollection: b5,
		ActiveCollectionID:  m["ActiveCollectionID"],
		HasActiveCart:       b6,
		IsViewingCart:       b7,
		IsInCheckout:        b8,
	}
}

func AddActive(rdb *redis.Client, tenant, dim, id, sessionID string) {
	now := time.Now()
	b := now.Unix() / 15
	k := fmt.Sprintf("active:%s:%s:%s:%d", tenant, dim, id, b)
	rdb.SAdd(context.TODO(), k, sessionID)
	rdb.Expire(context.TODO(), k, 4*time.Minute+15*time.Second)
}
