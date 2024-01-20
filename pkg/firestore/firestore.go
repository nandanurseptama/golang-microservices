package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type FirestoreClient = firestore.Client

// Initiate firestore client
//
// return client if success
//
// return error if failed
func NewClient(serviceAccountPath string) (*FirestoreClient, error) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile(serviceAccountPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		logrus.Error(err)
	}

	return app.Firestore(ctx)
}
