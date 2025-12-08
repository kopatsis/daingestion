package live

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func AddActive(rdb *redis.Client, tenant, dim, id, sessionID string) {
	now := time.Now()
	b := now.Unix() / 10
	k := fmt.Sprintf("active:%s:%s:%s:%d", tenant, dim, id, b)
	rdb.SAdd(context.TODO(), k, sessionID)
	rdb.Expire(context.TODO(), k, 4*time.Minute)
}
