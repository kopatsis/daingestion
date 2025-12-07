package live

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func AddActive(rdb *redis.Client, tenant, dim, id, sessionID string) {
	now := time.Now()
	b := now.Unix() / 10
	k := fmt.Sprintf("active:%s:%s:%s:%d", tenant, dim, id, b)
	rdb.SAdd(k, sessionID)
	rdb.Expire(k, 4*time.Minute)
}

func AddActiveMap(rdb *redis.Client, tenant, dim, id, sessionID string) {
	now := time.Now()
	b := now.Unix() / 20
	k := fmt.Sprintf("active:%s:%s:%s:%d", tenant, dim, id, b)
	rdb.SAdd(k, sessionID)
	rdb.Expire(k, 4*time.Minute)
}
