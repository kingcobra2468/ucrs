package device

import (
	"time"

	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   DeviceService
}

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
