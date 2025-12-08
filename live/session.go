package live

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	SessionNewClient = "new_client"
	SessionExpired   = "expired_session"
	SessionValid     = "valid_session"
)

type Result struct {
	SessionID string
	Status    string
}

func ManageSession(ctx context.Context, rdb *redis.Client, clientID, store string) (Result, error) {
	key := "sess:" + store + ":" + clientID
	now := time.Now().Unix()

	vals, err := rdb.HMGet(ctx, key, "id", "ts").Result()
	if err != nil {
		return Result{}, err
	}

	idRaw := vals[0]
	tsRaw := vals[1]

	if idRaw == nil || tsRaw == nil {
		newID := "PXID-" + uuid.NewString()
		_, err = rdb.HSet(ctx, key, "id", newID, "ts", now).Result()
		if err != nil {
			return Result{}, err
		}
		return Result{SessionID: newID, Status: SessionNewClient}, nil
	}

	id := idRaw.(string)
	tsInt, _ := tsRaw.(int64)

	if now-tsInt > 1800 {
		newID := "PXID-" + uuid.NewString()
		_, err = rdb.HSet(ctx, key, "id", newID, "ts", now).Result()
		if err != nil {
			return Result{}, err
		}
		return Result{SessionID: newID, Status: SessionExpired}, nil
	}

	_, err = rdb.HSet(ctx, key, "ts", now).Result()
	if err != nil {
		return Result{}, err
	}

	return Result{SessionID: id, Status: SessionValid}, nil
}
