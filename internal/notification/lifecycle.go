package notification

import (
	"context"
	"errors"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

var client *messaging.Client

func ConnectFcm(ctx context.Context) error {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return err
	}

	client, err = app.Messaging(ctx)
	if err != nil {
		return err
	}

	return nil
}

func AddRT(ctx context.Context, token string) error {
	response, err := client.SubscribeToTopic(ctx, []string{token}, "un")
	if err != nil {
		return err
	}

	if response.FailureCount == 0 {
		return nil
	}

	return errors.New(response.Errors[0].Reason)
}

func RemoveRT(ctx context.Context, token string) error {
	response, err := client.UnsubscribeFromTopic(ctx, []string{token}, "un")
	if err != nil {
		return err
	}

	if response.FailureCount == 0 {
		return nil
	}

	return errors.New(response.Errors[0].Reason)
}
