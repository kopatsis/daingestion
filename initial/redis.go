package initial

import (
	"github.com/go-redis/redis"
)

func NewRedis() *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := r.Ping().Err(); err != nil {
		panic(err)
	}
	return r
}
