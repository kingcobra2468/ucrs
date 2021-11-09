package device

import (
	"context"
	"errors"

	"github.com/kingcobra2468/ucrs/internal/notification"
	"github.com/kingcobra2468/ucrs/internal/registry"
)

// Service representing ucrs microservice.
type DeviceService interface {
	RegisterToken(string) (bool, error)
	RefreshTokenTTL(string) (bool, error)
	UpdateToken(string, string) (bool, error)
}

// Device for managing clients to FCM and Redis.
type Device struct {
	Ds *notification.DeviceSubscriber
	Dr *registry.DatabaseRegistry
}

var (
	ErrTokenLifecycle = errors.New("unable to change lifecycle for the registration token")
)

// Add the registration token to a cache (which monitors for for stale tokens)
// and add it the "un" topic for recieving notifications.
func (d Device) RegisterToken(token string) (bool, error) {
	err := d.Dr.AddRegistrationToken(token)
	if err != nil {
		return false, ErrTokenLifecycle
	}
	err = d.Ds.AddRT(context.Background(), token)
	if err != nil {
		return false, ErrTokenLifecycle
	}

	return true, nil
}

// Reset the TTL decay of a registration token to its initial value.
func (d Device) RefreshTokenTTL(token string) (bool, error) {
	status, err := d.Dr.ResetTokenTTL(token)
	if err != nil || !status {
		return false, ErrTokenLifecycle
	}

	return true, nil
}

// Removes the old registration token from the registry and FCM topic. Then
// adds the new token to both the registry and the FCM topic.
func (d Device) UpdateToken(newToken, oldToken string) (bool, error) {
	status, err := d.Dr.UpdateToken(newToken, oldToken)
	if err != nil || !status {
		return false, ErrTokenLifecycle
	}

	err = d.Ds.RemoveRT(context.Background(), oldToken)
	if err != nil {
		return false, ErrTokenLifecycle
	}

	err = d.Ds.AddRT(context.Background(), newToken)
	if err != nil {
		return false, ErrTokenLifecycle
	}

	return true, nil
}
