package firebase

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Client struct {
	authClient *auth.Client
}

func NewClient(ctx context.Context) (*Client, error) {
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
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{authClient: authClient}, nil
}

func (c *Client) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return c.authClient.VerifyIDToken(ctx, idToken)
}

func (c *Client) CreateUser(ctx context.Context, email, password string) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	return c.authClient.CreateUser(ctx, params)
}

// AuthProvider adapts Client to the user.AuthProvider interface.
type AuthProvider struct {
	client *Client
}

func NewAuthProvider(client *Client) *AuthProvider {
	return &AuthProvider{client: client}
}

func (p *AuthProvider) CreateUser(ctx context.Context, email, password string) (string, error) {
	record, err := p.client.CreateUser(ctx, email, password)
	if err != nil {
		return "", err
	}
	return record.UID, nil
}
