package live

import (
	"context"
	"dmd/logging"

	"github.com/redis/go-redis/v9"
)

func CheckStoreAndIncrement(ctx context.Context, rdb *redis.Client, store, param, requestID string, isnew bool) (bool, error) {
	storeKey := "store:" + store

	exists, err := rdb.Exists(ctx, storeKey).Result()
	if err != nil {
		logging.LogError(
			"ERROR",
			"plan_check_failed",
			"http",
			store,
			param,
			requestID,
			false,
			"unable to confirm store exists",
		)
		return false, err
	}
	if exists == 0 {
		return false, nil
	}

	if isnew {
		counterKey := "sessions:" + store

		_, err = rdb.Incr(ctx, counterKey).Result()
		if err != nil {
			logging.LogError(
				"ERROR",
				"plan_check_failed",
				"http",
				store,
				param,
				requestID,
				false,
				"unable to increment session count",
			)

			return true, err
		}
	}

	return true, nil
}
