package device

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type AuthenticateRequest struct {
	Secret string `json:"secret"`
}

type AuthenticateResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

type RegisterTokenResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

func MakeAuthenticateEndpoint(ds DeviceService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthenticateRequest)
		fmt.Print(req)

		return AuthenticateResponse{true, ErrAuthInvalid}, nil
	}
}
