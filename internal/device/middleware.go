package device

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Middlewere for performing request-based logging of the endpoints.
type LoggingMiddleware struct {
	Logger log.Logger
	Next   DeviceService
}

// Logging wrapper for token registration logic.
func (lm LoggingMiddleware) RegisterToken(token string) (bool, error) {
	defer func(begin time.Time) {
		lm.Logger.Log(
			"method", "registertoken",
			"took", time.Since(begin),
		)
	}(time.Now())

	status, err := lm.Next.RegisterToken(token)
	return status, err
}

// Logging wrapper for token heartbeat check.
func (lm LoggingMiddleware) RefreshTokenTTL(token string) (bool, error) {
	defer func(begin time.Time) {
		lm.Logger.Log(
			"method", "refreshtokenttl",
			"took", time.Since(begin),
		)
	}(time.Now())

	status, err := lm.Next.RefreshTokenTTL(token)
	return status, err
}

// Logging wrapper for token expiration reset logic.
func (lm LoggingMiddleware) UpdateToken(newToken, oldToken string) (bool, error) {
	defer func(begin time.Time) {
		lm.Logger.Log(
			"method", "updatetoken",
			"took", time.Since(begin),
		)
	}(time.Now())

	status, err := lm.Next.UpdateToken(newToken, oldToken)
	return status, err
}
