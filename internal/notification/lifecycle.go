package notification

import (
	"context"
	"errors"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type Subscription interface {
	Connect(ctx context.Context) error
	AddRT(ctx context.Context, token string) error
	RemoveRT(ctx context.Context, token string) error
}

type DeviceSubscriber struct {
	client *messaging.Client
}

var (
	ErrNotConnected = errors.New("client not connected")
)

func (ds *DeviceSubscriber) Connect(ctx context.Context) error {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return err
	}

	ds.client, err = app.Messaging(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DeviceSubscriber) AddRT(ctx context.Context, token string) error {
	if ds.client == nil {
		return ErrNotConnected
	}

	response, err := ds.client.SubscribeToTopic(ctx, []string{token}, "un")
	if err != nil {
		return err
	}

	if response.FailureCount == 0 {
		return nil
	}

	return errors.New(response.Errors[0].Reason)
}

func (ds *DeviceSubscriber) RemoveRT(ctx context.Context, token string) error {
	if ds.client == nil {
		return ErrNotConnected
	}

	response, err := ds.client.UnsubscribeFromTopic(ctx, []string{token}, "un")
	if err != nil {
		return err
	}

	if response.FailureCount == 0 {
		return nil
	}

	return errors.New(response.Errors[0].Reason)
}
