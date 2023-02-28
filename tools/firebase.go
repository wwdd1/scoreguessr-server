package tools

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var client *auth.Client

func GetFirebaseClient(ctx context.Context) (*auth.Client, error) {
	if client != nil {
		return client, nil
	}

	var err error
	var app *firebase.App

	app, err = firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}

	client, err = app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
