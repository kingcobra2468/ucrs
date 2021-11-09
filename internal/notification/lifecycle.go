package notification

import (
	"context"
	"errors"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

// Handle for Google FCM. Managing the lifecycle of
// adding and removing registration tokens from the specified topic
// that is used for broadcasting notifcations.
type Subscription interface {
	Connect(ctx context.Context) error
	AddRT(ctx context.Context, token string) error
	RemoveRT(ctx context.Context, token string) error
}

// Client for registration token lifecycle.
type DeviceSubscriber struct {
	client *messaging.Client
	Topic  string
}

var (
	ErrNotConnected = errors.New("client not connected")
)

// Connect to Google's FCM.
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

// Add registration token to the specified topic.
func (ds *DeviceSubscriber) AddRT(ctx context.Context, token string) error {
	if ds.client == nil {
		return ErrNotConnected
	}

	response, err := ds.client.SubscribeToTopic(ctx, []string{token}, ds.Topic)
	if err != nil {
		return err
	}

	// Check if FCM didn't return any failures which indicates succcess.
	if response.FailureCount == 0 {
		return nil
	}

	return errors.New(response.Errors[0].Reason)
}

// Remove registration token to the specified topic.
func (ds *DeviceSubscriber) RemoveRT(ctx context.Context, token string) error {
	if ds.client == nil {
		return ErrNotConnected
	}

	response, err := ds.client.UnsubscribeFromTopic(ctx, []string{token}, ds.Topic)
	if err != nil {
		return err
	}
	// Check if FCM didn't return any failures which indicates succcess.
	if response.FailureCount == 0 {
		return nil
	}

	return errors.New(response.Errors[0].Reason)
}
