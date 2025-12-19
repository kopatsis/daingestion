package live

import (
	"context"
	"dmd/models"
	"strings"

	"github.com/redis/go-redis/v9"
)

func MainLiveWork(client *redis.Client, sessionStruct SessionActiveState, eventID, clientID, store, url, param string, ev *models.IngestEvent) (bool, error) {

	isDedupEvent, err := DedupEventID(client, store, eventID)
	if err != nil {
		return false, err
	} else if isDedupEvent {
		return true, nil
	}

	isDedupView, err := Dedup(client, store, clientID, url, param)
	if err != nil {
		return false, err
	} else if isDedupView {
		return true, nil
	}

	sessionResults, err := ManageSession(context.Background(), client, clientID, store)
	if err != nil {
		return false, err
	}

	if err := SetState(client, sessionResults.SessionID, sessionStruct); err != nil {
		return false, err
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
		return false, err
	}

	return false, nil
}
