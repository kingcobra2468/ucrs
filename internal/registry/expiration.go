package registry

import (
	"errors"
)

// Signature for creating callbacks for token expiration.
type OnTokenExpiration func(token string) error

// Listener for handling token expiration.
func (dr *DatabaseRegistry) NewTokenExpirationListener(done <-chan bool, ote OnTokenExpiration) error {
	pubsub := dr.rdb.PSubscribe("__keyevent@0__:expired")
	_, err := pubsub.Receive()
	if err != nil {
		return errors.New("unable to start token expiration listener")
	}

	ch := pubsub.Channel()

	// Launch listener in a background thread.
	go func() {
		// Operate on expired tokens without blocking.
		go func() {
			for msg := range ch {
				ote(msg.Payload)
			}
		}()
		// Operate listener until terminal.
		<-done
		pubsub.Close()
	}()

	return nil
}
