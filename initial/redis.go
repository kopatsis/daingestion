package initial

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := r.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return r
}
