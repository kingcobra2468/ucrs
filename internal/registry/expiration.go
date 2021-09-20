package registry

import (
	"errors"
)

type OnTokenExpiration func(token string)

func StartTokenExpirationListener(done <-chan bool, ote OnTokenExpiration) error {
	pubsub := cache.rdb.PSubscribe("__keyevent@0__:expired")
	_, err := pubsub.Receive()
	if err != nil {
		return errors.New("unable to start token expiration listener")
	}

	ch := pubsub.Channel()

	// launch listener in a background thread
	go func() {
		// operate on expired tokens without blocking
		go func() {
			for msg := range ch {
				// fmt.Println(msg.Channel, msg.Payload)
				ote(msg.Payload)
			}
		}()

		<-done
		pubsub.Close()
	}()

	return nil
}
