package lib

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Firebase struct {
	App     *firebase.App
	Context context.Context
	Client  *messaging.Client
	Logger  Logger
}

func NewFirebase(env Env, logger Logger) Firebase {

	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs(env.FirebaseConfigPath)
	if err != nil {
		logger.Error("Unable to locate firebase config file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Error("Firebase SDK initialization failed")
	}

	//Messaging client
	client, _ := app.Messaging(ctx)

	return Firebase{
		App:     app,
		Context: ctx,
		Client:  client,
		Logger:  logger,
	}
}

func (f Firebase) SendToToken(registrationToken string, title string, body string, data map[string]string) error {
	ctx := context.Background()
	client, err := f.App.Messaging(ctx)
	if err != nil {
		f.Logger.Errorf("error getting Messaging client: %v\n", err)
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: registrationToken,
		Data:  data,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		f.Logger.Errorf("error sending message: %v\n", err)
	}
	f.Logger.Info("Successfully sent message:", response)
	return err
}
