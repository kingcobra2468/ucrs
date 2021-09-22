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

type Device struct{}

var (
	ErrAuthInvalid = errors.New("failued to authenticate")
)

func (d Device) Authenticate(secret string) (bool, error) {
	return true, nil
}

func (d Device) RegisterToken(token string) (bool, error) {
	registry.AddRegistrationToken(token)
	notification.AddRT(context.Background(), token)

	return true, nil
}
