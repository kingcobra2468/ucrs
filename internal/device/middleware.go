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

// Logging wrapper for authentication logic.
func (lm LoggingMiddleware) Authenticate(secret string) (bool, error) {
	defer func(begin time.Time) {
		lm.Logger.Log(
			"method", "authenticate",
			"took", time.Since(begin),
		)
	}(time.Now())

	status, err := lm.Next.Authenticate(secret)
	return status, err
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
