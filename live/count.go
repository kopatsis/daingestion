package live

import (
	"context"
	"dmd/bots"
	"dmd/models"
	"dmd/steps"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func CreateSessionStruct(ev models.IngestEvent, geo models.GeoData, uaInfo models.UAInfo, utm models.UTM, pageType steps.PageType, botScore bots.BotLevel, ref models.Referrer, param string) models.SessionActiveState {
	sessionStruct := models.SessionActiveState{
		Country:     geo.CountryISO,
		Region:      geo.SubdivisionISO,
		City:        geo.CityName,
		Latitude:    geo.Latitude,
		Longitude:   geo.Longitude,
		IsBot:       botScore > 0,
		ASN:         geo.ASN,
		ASNProvider: geo.DataCenter,
		DeviceType:  uaInfo.Type,
		DeviceBrand: uaInfo.Brand,
		OSName:      uaInfo.OSName,
		BrowserName: uaInfo.ClientName,
		Language:    ev.Event.Context.Navigator.Language,
		RefDomain:   ref.DomainOnly,
		UTMSource:   utm.Source,
		RouteType:   string(pageType),
		Route:       ev.Event.Context.Document.Location.Pathname,
		FullURL:     ev.Event.Context.Document.Location.Href,
	}

	if ev.Init.Data.Customer != nil {
		sessionStruct.IsLoggedIn = true
		sessionStruct.CustomerID = ev.Init.Data.Customer.ID
		sessionStruct.HasPreviousPurchase = ev.Init.Data.Customer.OrdersCount > 0
	}

	if param == "product_viewed" {
		variantID, productID, err := ExtractProductIDs(ev.Event.Data)
		sessionStruct.IsViewingProduct = true
		if err == nil {
			sessionStruct.ActiveProductID = productID
			sessionStruct.ActiveVariantID = variantID
		}
	} else if param == "collection_viewed" {
		collectionID, err := ExtractCollectionID(ev.Event.Data)
		sessionStruct.IsViewingCollection = true
		if err == nil {
			sessionStruct.ActiveCollectionID = collectionID
		}
	} else if param == "cart_viewed" {
		sessionStruct.IsViewingCart = true
	} else if strings.Contains(param, "checkout") && param != "checkout_completed" {
		sessionStruct.IsInCheckout = true
	}

	if models.CheckIfCart(ev.Event.Data) {
		sessionStruct.HasActiveCart = true
	}

	return sessionStruct
}

func SetState(rdb *redis.Client, sessionID string, s models.SessionActiveState) error {
	key := "session:" + sessionID
	_, err := rdb.HSet(context.TODO(), key, map[string]interface{}{
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
	return rdb.Expire(context.TODO(), key, 300000000000).Err()
}

func GetState(ctx context.Context, rdb *redis.Client, sessionID string) models.SessionActiveState {
	key := "session:" + sessionID
	m, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return models.SessionActiveState{}
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

	return models.SessionActiveState{
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

func AddActive(rdb *redis.Client, store, sessionID string) error {
	now := time.Now()
	b := now.Unix() / 15
	k := fmt.Sprintf("active:%s:%d", store, b)

	if err := rdb.SAdd(context.TODO(), k, sessionID).Err(); err != nil {
		return err
	}

	if err := rdb.Expire(context.TODO(), k, 4*time.Minute+15*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func GetActiveLast16Buckets(rdb *redis.Client, store string) ([]string, error) {
	now := time.Now().Unix()
	currentBucket := now / 15

	if now%15 != 0 {
		currentBucket--
	}

	startBucket := currentBucket - 15

	ctx := context.TODO()
	combined := make(map[string]struct{})

	for b := startBucket; b <= currentBucket; b++ {
		key := fmt.Sprintf("active:%s:%d", store, b)
		members, err := rdb.SMembers(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		for _, m := range members {
			combined[m] = struct{}{}
		}
	}

	out := make([]string, 0, len(combined))
	for k := range combined {
		out = append(out, k)
	}

	return out, nil
}
