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

func loadCredentialsJSON() ([]byte, error) {
	encoded := os.Getenv("FIREBASE_CREDENTIALS")
	if encoded != "" {
		credentialsJSON, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, errors.New("failed to decode FIREBASE_CREDENTIALS: must be base64 encoded")
		}
		return credentialsJSON, nil
	}

	candidates := []string{}
	if path := os.Getenv("FIREBASE_CREDENTIALS_FILE"); path != "" {
		candidates = append(candidates, path)
	}
	candidates = append(candidates, "firebase-credentials.json", "../../firebase-credentials.json")

	for _, path := range candidates {
		credentialsJSON, err := os.ReadFile(path)
		if err == nil {
			return credentialsJSON, nil
		}
	}

	return nil, errors.New("firebase credentials not found: set FIREBASE_CREDENTIALS env or provide firebase-credentials.json (or FIREBASE_CREDENTIALS_FILE)")
}

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	credentialsJSON, err := loadCredentialsJSON()
	if err != nil {
		return nil, err
	}

	opt := option.WithAuthCredentialsJSON(option.ServiceAccount, credentialsJSON)
	client, err := firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		return nil, err
	}
	return &Client{Client: client}, nil
}
