package registry

import (
	"time"

	"github.com/go-redis/redis"
)

type Registry interface {
	Connect(addr string)
	AddRegistrationToken(token string) error
	StartTokenExpirationListener(done <-chan bool, ote OnTokenExpiration) error
}

type DatabaseRegistry struct {
	rdb *redis.Client
}

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

func (dr *DatabaseRegistry) AddRegistrationToken(token string) error {
	err := dr.rdb.Set(token, token, time.Hour*48).Err()

	return err
}
