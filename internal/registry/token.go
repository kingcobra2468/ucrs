package registry

import (
	"time"

	"github.com/go-redis/redis"
)

// Cache for handling lifecycle/staleness of registration tokens.
type Registry interface {
	Connect(addr string)
	AddRegistrationToken(token string) error
	StartTokenExpirationListener(done <-chan bool, ote OnTokenExpiration) error
}

// Client for cache.
type DatabaseRegistry struct {
	rdb *redis.Client
}

// Connect to Redis which will be used for caching.
func (dr *DatabaseRegistry) Connect(addr string) {
	if dr.rdb != nil {
		return
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	dr.rdb = client
}

// Add a registration token to the cache with the given TTL until
// expiration.
func (dr *DatabaseRegistry) AddRegistrationToken(token string) error {
	err := dr.rdb.Set(token, token, time.Hour*48).Err()

	return err
}
