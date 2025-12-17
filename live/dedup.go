package live

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func Dedup(client *redis.Client, store, clientID, url, param string) (bool, error) {

	now := time.Now().Unix()

	if !(strings.Contains(param, "viewed") || strings.Contains(param, "started") || param == "search_submitted") {
		return true, nil
	}

	key := "dedupe:" + store + ":" + clientID
	v, err := client.Get(context.TODO(), key).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	if err == nil {
		var prev struct {
			URL string
			TS  int64
		}
		_ = json.Unmarshal([]byte(v), &prev)
		if prev.URL == url && now-prev.TS < 4 {
			return false, nil
		}
	}
	b, _ := json.Marshal(struct {
		URL string
		TS  int64
	}{URL: url, TS: now})
	err = client.Set(context.TODO(), key, b, 5*time.Second).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func DedupEventID(client *redis.Client, store, eventID string) (bool, error) {
	key := "eventid:" + store + ":" + eventID
	ok, err := client.SetNX(context.TODO(), key, "", 90*time.Second).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}
