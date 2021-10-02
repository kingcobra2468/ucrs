package device

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

// Schema of JSON object for authentication request
type AuthenticateRequest struct {
	Secret string `json:"secret"`
}

// Schema of JSON object of token registration request
type RegisterTokenRequest struct {
	RegistrationToken string `json:"registration_token"`
}

// Schema of JSON object of token alive request
type RefreshTokenTTLRequest struct {
	RegistrationToken string `json:"registration_token"`
}

// Schema of JSON object of token update request
type UpdateTokenRequest struct {
	OldToken string `json:"old_token"`
	NewToken string `json:"new_token"`
}

// Schema of JSON object for authentication response
type AuthenticateResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

// Schema of JSON object of token registration response
type RegisterTokenResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

// Schema of JSON object of token registration response
type RefreshTokenTTLResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

// Schema of JSON object of token registration response
type UpdateTokenResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

// Endpoint for authentication. Currently this is dummy endpoint as
// it contains no logic for authentication.
func makeAuthenticateEndpoint(ds DeviceService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthenticateRequest)
		fmt.Print(req)

		return AuthenticateResponse{true, ErrAuthInvalid}, nil
	}
}

// Endpoint for token registration. Cache the token and add it to
// the "un" topic.
func makeRegisterTokenEndpoint(ds DeviceService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterTokenRequest)
		ds.RegisterToken(req.RegistrationToken)

		return RegisterTokenResponse{true, ErrAuthInvalid}, nil
	}
}

func makeRefreshTokenTTLEndpoint(ds DeviceService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(RefreshTokenTTLRequest)
		ds.RefreshTokenTTL(req.RegistrationToken)

		return RegisterTokenResponse{true, ErrAuthInvalid}, nil
	}
}

func makeUpdateTokenEndpoint(ds DeviceService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateTokenRequest)
		ds.UpdateToken(req.NewToken, req.OldToken)

		return RegisterTokenResponse{true, ErrAuthInvalid}, nil
	}
}
