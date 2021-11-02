package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/go-kit/kit/log"
	"github.com/kingcobra2468/ucrs/internal/device"
	"github.com/kingcobra2468/ucrs/internal/notification"
	"github.com/kingcobra2468/ucrs/internal/registry"
)

type config struct {
	ServiceHostname string `env:"UCRS_HOSTNAME" envDefault:"127.0.0.1"`
	ServicePort     int    `env:"UCRS_PORT" envDefault:"8080"`
	redisHostname   string `env:"UCRS_REDIS_HOSTNAME" envDefault:"127.0.0.1"`
	redisPort       int    `env:"UCRS_REDIS_PORT" envDefault:"8080"`
	FcmTopic        string `env:"UCRS_FCM_TOPIC" envDefault:"un"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var ds notification.DeviceSubscriber = notification.DeviceSubscriber{Topic: cfg.FcmTopic}
	{
		ds.Connect(context.Background())
	}

	var dr registry.DatabaseRegistry = registry.DatabaseRegistry{}
	{
		dr.Connect(fmt.Sprintf("%s:%d", cfg.redisHostname, cfg.redisPort))
	}

	done := make(chan bool)

	dr.NewTokenExpirationListener(done, func(token string) error {
		fmt.Printf("UnSubscribing Token: %s\n", token)
		ds.RemoveRT(context.Background(), token)

		return nil
	})

	var service device.DeviceService = device.Device{Ds: &ds, Dr: &dr}
	service = device.LoggingMiddleware{Logger: logger, Next: service}
	var h http.Handler = device.MakeHTTPHandler(service)

	errs := make(chan error)
	// Listener for Ctrl+C signals.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	// Launch microservice.
	go func() {
		url := fmt.Sprintf("%s:%d", cfg.ServiceHostname, cfg.ServicePort)

		logger.Log("transport", "HTTP", "addr", url)
		errs <- http.ListenAndServe(url, h)
	}()

	logger.Log("exit", <-errs)
}
