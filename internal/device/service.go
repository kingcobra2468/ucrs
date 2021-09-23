package device

import (
	"context"
	"errors"

	"github.com/kingcobra2468/ucrs/internal/notification"
	"github.com/kingcobra2468/ucrs/internal/registry"
)

type DeviceService interface {
	Authenticate(string) (bool, error)
	RegisterToken(string) (bool, error)
}

type Device struct {
	Ds *notification.DeviceSubscriber
	Dr *registry.DatabaseRegistry
}

var (
	ErrAuthInvalid = errors.New("failed to authenticate")
)

func (d Device) Authenticate(secret string) (bool, error) {
	return true, nil
}

func (d Device) RegisterToken(token string) (bool, error) {
	d.Dr.AddRegistrationToken(token)
	d.Ds.AddRT(context.Background(), token)

	return true, nil
}
