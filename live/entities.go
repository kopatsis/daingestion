package live

import (
	"context"
	"dmd/models"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func AllLiveUpdates(ev *models.IngestEvent, rdb *redis.Client, store, eventType string) error {

	if err := UpdateCustomer(ev, rdb, store); err != nil {
		return err
	}

	return nil
}

type customerPayload struct {
	Customer *struct {
		ID string `json:"id"`
	} `json:"customer"`
}

func UpdateCustomer(ev *models.IngestEvent, rdb *redis.Client, store string) error {
	if len(ev.Init.Data.Rest) == 0 {
		return nil
	}

	var p customerPayload
	if err := json.Unmarshal(ev.Init.Data.Rest, &p); err != nil {
		return nil
	}

	if p.Customer == nil {
		return nil
	}
	if p.Customer.ID == "" {
		return nil
	}

	key := "lastseen:" + store + ":cust:" + p.Customer.ID
	_, err := rdb.Set(context.TODO(), key, time.Now().Unix(), 0).Result()
	return err
}
