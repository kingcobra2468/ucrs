package device

import (
	"context"
	"errors"

	"github.com/kingcobra2468/ucrs/internal/notification"
	"github.com/kingcobra2468/ucrs/internal/registry"
)

// Service representing ucrs microservice.
type DeviceService interface {
	Authenticate(string) (bool, error)
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
	ErrAuthInvalid = errors.New("failed to authenticate")
)

// Authenticate device prior to adding its registration token to
// the FCM.
func (d Device) Authenticate(secret string) (bool, error) {
	return true, nil
}

// Add the registration token to a cache (which monitors for for stale tokens)
// and add it the "un" topic for recieving notifications.
func (d Device) RegisterToken(token string) (bool, error) {
	d.Dr.AddRegistrationToken(token)
	d.Ds.AddRT(context.Background(), token)

	return true, nil
}

func (d Device) RefreshTokenTTL(token string) (bool, error) {
	return true, nil
}

func (d Device) UpdateToken(newToken, oldToken string) (bool, error) {
	return true, nil
}
