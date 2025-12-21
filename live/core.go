package live

import (
	"context"
	"dmd/models"
	"strings"

	"github.com/redis/go-redis/v9"
)

func MainLiveWork(client *redis.Client, sessionStruct models.SessionActiveState, eventID, clientID, store, url, param, reqID string, ev *models.IngestEvent) (Result, bool, error) {

	isDedupEvent, err := DedupEventID(client, store, param, eventID, reqID)
	if err != nil {
		return Result{}, false, err
	} else if isDedupEvent {
		return Result{}, true, nil
	}

	isDedupView, err := Dedup(client, store, clientID, url, param, reqID)
	if err != nil {
		return Result{}, false, err
	} else if isDedupView {
		return Result{}, true, nil
	}

	sessionResults, err := ManageSession(context.Background(), client, clientID, store, param, reqID)
	if err != nil {
		return Result{}, false, err
	}

	allowedStore, err := CheckStoreAndIncrement(context.Background(), client, store, param, reqID, sessionResults.Status == SessionNewClient)
	if err != nil {
		return Result{}, false, err
	} else if !allowedStore {
		return Result{}, true, nil
	}

	if err := SetState(client, sessionResults.SessionID, reqID, store, param, sessionStruct); err != nil {
		return sessionResults, false, err
	}

	location := ""
	if sessionStruct.Country == "" {
		location = "Unknown Location"
	} else {
		if sessionStruct.City != "" {
			location = sessionStruct.City + ", "
		}
		if sessionStruct.Region != "" {
			location += sessionStruct.Region + ", "
		}
		location += sessionStruct.Country
	}

	lv := LiveEvent{
		Store:     store,
		EventCode: param,
		SessionID: sessionResults.SessionID,
		Device:    sessionStruct.DeviceType,
		Location:  location,
		URL:       url,
	}.WithHumanizedEvent()

	if param == "search_submitted" {
		searchQuery, err := ExtractSearchQuery(ev.Event.Data)
		if err == nil {
			lv.Search = &searchQuery
		}
	} else if param == "checkout_completed" {
		orderID, err := ExtractCheckoutOrderID(ev.Event.Data)
		if err == nil {
			lv.OrderID = &orderID
		}
	} else if param == "collection_viewed" {
		collectionTitle, err := ExtractCollectionTitle(ev.Event.Data)
		if err == nil {
			lv.Collection = &collectionTitle
		}
	} else if strings.Contains(param, "product") {
		orderID, err := ExtractCheckoutOrderID(ev.Event.Data)
		if err == nil {
			lv.OrderID = &orderID
		}
	}

	if err := PublishEvent(context.TODO(), client, lv); err != nil {
		return sessionResults, false, err
	}

	return sessionResults, false, nil
}
