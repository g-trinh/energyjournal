package firestore

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type Client struct {
	*firestore.Client
}

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	credentialsJSON, err := os.ReadFile("../../firebase-credentials.json")
	if err != nil {
		// Try from env
		encoded := os.Getenv("FIREBASE_CREDENTIALS")
		if encoded == "" {
			return nil, errors.New("firebase credentials not found: set FIREBASE_CREDENTIALS env or provide firebase-credentials.json")
		}
		credentialsJSON, err = base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, errors.New("failed to decode FIREBASE_CREDENTIALS: must be base64 encoded")
		}
	}

	opt := option.WithAuthCredentialsJSON(option.ServiceAccount, credentialsJSON)
	client, err := firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		return nil, err
	}
	return &Client{Client: client}, nil
}
