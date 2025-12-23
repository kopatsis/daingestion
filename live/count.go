package live

import (
	"context"
	"dmd/bots"
	"dmd/logging"
	"dmd/models"
	"dmd/steps"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func CreateSessionStruct(ev models.IngestEvent, geo models.GeoData, uaInfo models.UAInfo, utm models.UTM, pageType steps.PageType, botScore bots.BotLevel, ref models.Referrer, store, param, reqID string) models.SessionActiveState {
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
		variantID, productID, err := ExtractProductIDs(ev.Event.Data, store, param, reqID)
		sessionStruct.IsViewingProduct = true
		if err == nil {
			sessionStruct.ActiveProductID = productID
			sessionStruct.ActiveVariantID = variantID
		}
	} else if param == "collection_viewed" {
		collectionID, err := ExtractCollectionID(ev.Event.Data, store, param, reqID)
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

func SetState(rdb *redis.Client, sessionID, requestID, store, param string, s models.SessionActiveState) error {
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
		logging.LogError(
			"FAILURE",
			"session_metadata_save_failed",
			"redis",
			store,
			param,
			requestID,
			true,
			"failed to persist session metadata",
		)

		return err
	}
	return rdb.Expire(context.TODO(), key, 300000000000).Err()
}

func AddActive(rdb *redis.Client, store, sessionID, param, requestID string) error {
	now := time.Now()
	b := now.Unix() / 15
	k := fmt.Sprintf("active:%s:%d", store, b)

	if err := rdb.SAdd(context.TODO(), k, sessionID).Err(); err != nil {
		logging.LogError(
			"FAILURE",
			"add_active_id_fail",
			"redis",
			store,
			param,
			requestID,
			false,
			"failed to add session ID to active list",
		)
		return err
	}

	if err := rdb.Expire(context.TODO(), k, 4*time.Minute+15*time.Second).Err(); err != nil {
		logging.LogError(
			"FAILURE",
			"add_active_id_fail",
			"redis",
			store,
			param,
			requestID,
			false,
			"failed to add expire session ID active list",
		)
		return err
	}

	return nil
}
