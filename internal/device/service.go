package device

import "errors"

type DeviceService interface {
	Authenticate(string) (bool, error)
	RegisterToken(string) (bool, error)
}

type Device struct{}

var (
	ErrAuthInvalid = errors.New("failued to authenticate")
)

func (d *Device) Authenticate(secret string) (bool, error) {
	return true, nil
}

func (d *Device) RegisterToken(secret string) (bool, error) {
	return true, nil
}
