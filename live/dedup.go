package live

import (
	"context"
	"dmd/logging"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func Dedup(client *redis.Client, store, clientID, url, param, requestID string) (bool, error) {

	if !(strings.Contains(param, "viewed") || strings.Contains(param, "started") || param == "search_submitted") {
		return true, nil
	}

	now := time.Now().Unix()

	key := "dedupe:" + store + ":" + clientID
	v, err := client.Get(context.TODO(), key).Result()
	if err != nil && err != redis.Nil {
		logging.LogError(
			"ERROR",
			"dedup_view_check_failed",
			"redis",
			store,
			param,
			requestID,
			true,
			"redis error viewing key during dedup check",
		)

		return false, err
	}
	if err == nil {
		var prev struct {
			URL string
			TS  int64
		}
		_ = json.Unmarshal([]byte(v), &prev)
		if prev.URL == url && now-prev.TS < 4 {
			return true, nil
		}
	}
	b, _ := json.Marshal(struct {
		URL string
		TS  int64
	}{URL: url, TS: now})
	err = client.Set(context.TODO(), key, b, 5*time.Second).Err()
	if err != nil {
		logging.LogError(
			"ERROR",
			"dedup_view_check_failed",
			"redis",
			store,
			param,
			requestID,
			true,
			"redis error setting key during dedup check",
		)
		return false, err
	}
	return false, nil
}

func DedupEventID(client *redis.Client, store, eventType, eventID, requestID string) (bool, error) {
	key := "eventid:" + store + ":" + eventID
	ok, err := client.SetNX(context.TODO(), key, "", 90*time.Second).Result()
	if err != nil {
		logging.LogError(
			"ERROR",
			"dedup_event_check_failed",
			"redis",
			store,
			eventType,
			requestID,
			true,
			"redis error during dedup check",
		)

		return false, err
	}
	return ok, nil
}
