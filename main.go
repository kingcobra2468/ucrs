package main

import (
	"context"
	"flag"
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
	hostname := flag.String("hostname", "0.0.0.0", "hostname for ucrs")
	port := flag.Int("port", 8080, "port for ucrs")
	redisHostname := flag.String("redis-hostname", "0.0.0.0", "hostname for redis cache")
	redisPort := flag.Int("redis-port", 6379, "port for redis cache")
	topic := flag.String("topic", "un", "fcm topic for registration token subscription")

	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var ds notification.DeviceSubscriber = notification.DeviceSubscriber{Topic: *topic}
	{
		ds.Connect(context.Background())
	}

	var dr registry.DatabaseRegistry = registry.DatabaseRegistry{}
	{
		dr.Connect(fmt.Sprintf("%s:%d", *redisHostname, *redisPort))
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
		url := fmt.Sprintf("%s:%d", *hostname, *port)

		logger.Log("transport", "HTTP", "addr", url)
		errs <- http.ListenAndServe(url, h)
	}()

	logger.Log("exit", <-errs)
}
