package registry

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// Cache for handling lifecycle/staleness of registration tokens.
type Registry interface {
	Connect(addr string)
	AddRegistrationToken(token string) error
	StartTokenExpirationListener(done <-chan bool, ote OnTokenExpiration) error
	ResetTokenTTL(token string) (bool, error)
	UpdateToken(newToken, oldToken string) (bool, error)
}

// Client for cache.
type DatabaseRegistry struct {
	rdb *redis.Client
}

var (
	ErrInvalidToken = errors.New("registration token doesn't exist")
	expirationDelta = time.Hour * 48
)

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
	err := dr.rdb.Set(token, token, expirationDelta).Err()

	return err
}

// Resets the TTL of a decaying registration token to the original decay
// value. Evaluates the status of whether the decay reset was successful.
func (dr *DatabaseRegistry) ResetTokenTTL(token string) (bool, error) {
	status, err := dr.rdb.Expire(token, expirationDelta).Result()
	if err != nil {
		return false, err
	}
	if !status {
		return false, ErrInvalidToken
	}

	return true, nil
}

// Updates an old registration token to a new value. Occurs when FCM clients
// recieve new token and thus must notify of this change to the server. Evaluates
// the status on whether the token update took place successfully.
func (dr *DatabaseRegistry) UpdateToken(newToken, oldToken string) (bool, error) {
	status, err := dr.rdb.Del(oldToken).Result()
	if err != nil {
		return false, err
	}
	if status == 0 {
		return false, ErrInvalidToken
	}

	err = dr.rdb.Set(newToken, newToken, expirationDelta).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}
