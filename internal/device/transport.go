package device

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type errorer interface {
	error() error
}

func MakeHTTPHandler(ds DeviceService) http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/api/auth").Handler(httptransport.NewServer(
		MakeAuthenticateEndpoint(ds),
		decodeAuthenticateRequest,
		encodeResponse,
	))

	r.Methods("POST").Path("/api/register-token").Handler(httptransport.NewServer(
		MakeRegisterTokenEndpoint(ds),
		decodeRegisterTokenRequest,
		encodeResponse,
	))

	return r
}

func decodeRegisterTokenRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req RegisterTokenRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeAuthenticateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req RegisterTokenRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrAuthInvalid:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
