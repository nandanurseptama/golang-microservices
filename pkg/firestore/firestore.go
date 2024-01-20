package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FirestoreClient struct {
	*firestore.Client
}

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
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return &FirestoreClient{client}, nil
}

func (f *FirestoreClient) GetCollection(path string) *CollectionRef {
	return &CollectionRef{f.Collection(path)}
}

type CollectionRef struct {
	*firestore.CollectionRef
}

func (c *CollectionRef) WhereColumn(path string, op string, value interface{}) Query {
	return Query{c.Where(path, op, value)}
}

type Query struct {
	firestore.Query
}

func (d Query) GetOne(
	ctx context.Context,
) (map[string]interface{}, error) {

	iter := d.Limit(1).Documents(ctx)

	var result map[string]interface{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		result = doc.Data()
	}
	return result, nil
}
