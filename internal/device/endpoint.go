package device

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

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

// Endpoint for token registration. Cache the token and add it to
// the specified topic.
func makeRegisterTokenEndpoint(ds DeviceService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterTokenRequest)
		_, err := ds.RegisterToken(req.RegistrationToken)

		if err != nil {
			return RegisterTokenResponse{false, err}, nil
		}

		return RegisterTokenResponse{true, nil}, nil
	}
}

// Endpoint for refreshing the decay TTL to its original decay value.
func makeRefreshTokenTTLEndpoint(ds DeviceService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(RefreshTokenTTLRequest)
		_, err := ds.RefreshTokenTTL(req.RegistrationToken)

		if err != nil {
			return RegisterTokenResponse{false, err}, nil
		}

		return RegisterTokenResponse{true, nil}, nil
	}
}

// Endpoint for updating an old token with its new counterpart given the FCM client.
func makeUpdateTokenEndpoint(ds DeviceService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateTokenRequest)
		_, err := ds.UpdateToken(req.NewToken, req.OldToken)

		if err != nil {
			return RegisterTokenResponse{false, err}, nil
		}

		return RegisterTokenResponse{true, nil}, nil
	}
}
