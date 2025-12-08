package live

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func CheckAndIncrement(ctx context.Context, rdb *redis.Client, planKey string, counterKey string) (bool, error) {
	planData, err := rdb.HGetAll(context.TODO(), planKey).Result()
	if err != nil {
		return false, err
	}
	limit, _ := strconv.ParseInt(planData["limit"], 10, 64)
	resetAt, _ := strconv.ParseInt(planData["reset_at"], 10, 64)
	now := time.Now().Unix()
	if now >= resetAt {
		t := time.Unix(resetAt, 0)
		for t.Unix() <= now {
			t = t.Add(30 * 24 * time.Hour)
		}
		newReset := t.Unix()
		ttl := newReset - now
		p := rdb.TxPipeline()
		p.HSet(context.TODO(), planKey, "reset_at", strconv.FormatInt(newReset, 10))
		p.Set(context.TODO(), counterKey, "1", time.Duration(ttl)*time.Second)
		_, err := p.Exec(context.TODO())
		if err != nil {
			return false, err
		}
		return false, nil
	}
	count, err := rdb.Incr(context.TODO(), counterKey).Result()
	if err != nil {
		return false, err
	}
	return count > limit, nil
}
