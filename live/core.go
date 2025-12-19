package live

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func MainLiveWork(client *redis.Client, sessionStruct SessionActiveState, eventID, clientID, store, url, param string) (bool, error) {

	isDedupEvent, err := DedupEventID(client, store, eventID)
	if err != nil {
		return false, err
	} else if isDedupEvent {
		return false, nil
	}

	isDedupView, err := Dedup(client, store, clientID, url, param)
	if err != nil {
		return false, err
	} else if isDedupView {
		return false, nil
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

	if err := PublishEvent(context.TODO(), client, LiveEvent{
		Store:     store,
		EventCode: param,
		SessionID: sessionResults.SessionID,
		Device:    sessionStruct.DeviceType,
		Location:  location,
	}.WithHumanizedEvent()); err != nil {
		return false, err
	}

	return true, nil
}
