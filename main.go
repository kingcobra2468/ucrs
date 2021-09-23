package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/kingcobra2468/ucrs/internal/device"
	"github.com/kingcobra2468/ucrs/internal/notification"
	"github.com/kingcobra2468/ucrs/internal/registry"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var ds notification.DeviceSubscriber = notification.DeviceSubscriber{}
	{
		ds.Connect(context.Background())
	}

	var dr registry.DatabaseRegistry = registry.DatabaseRegistry{}
	{
		dr.Connect("10.0.1.10:6389")
	}

	done := make(chan bool)

	dr.NewTokenExpirationListener(done, func(token string) error {
		fmt.Printf("UnSubscribing Token: %s\n", token)
		ds.RemoveRT(context.Background(), token)

		return nil
	})

	var service device.DeviceService = device.Device{Ds: &ds, Dr: &dr}

	var h http.Handler = device.MakeHTTPHandler(service)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", ":8080")
		errs <- http.ListenAndServe(":8080", h)
	}()

	logger.Log("exit", <-errs)
}
