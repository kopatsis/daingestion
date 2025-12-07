package live

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type Event struct {
	TenantID  string `json:"tenant_id"`
	EventType string `json:"event_type"`
	Message   string `json:"message"`
	Time      int64  `json:"time"`
}

func PublishEvent(ctx context.Context, rdb *redis.Client, e Event) error {
	ch := fmt.Sprintf("events:%s:%s", e.TenantID, e.EventType)
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return rdb.Publish(ch, b).Err()
}
