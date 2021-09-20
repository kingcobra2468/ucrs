package registry

import (
	"time"

	"github.com/go-redis/redis"
)

type registryTokenCache struct {
	rdb *redis.Client
}

var cache *registryTokenCache

func Connect(addr string) {
	if cache != nil {
		return
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	cache = &registryTokenCache{rdb: client}
}

func AddRegistrationToken(token string) error {
	err := cache.rdb.Set(token, token, time.Hour*48).Err()

	return err
}
