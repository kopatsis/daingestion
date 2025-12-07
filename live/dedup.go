package live

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

func Dedup(client *redis.Client, store, clientID, url string, now int64) (bool, error) {
	key := "dedupe:" + store + ":" + clientID
	v, err := client.Get(key).Result()
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
	err = client.Set(key, b, 5*time.Second).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func DedupEventID(client *redis.Client, store, eventID string) (bool, error) {
	key := "eventid:" + store + ":" + eventID
	ok, err := client.SetNX(key, "", 90*time.Second).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}
