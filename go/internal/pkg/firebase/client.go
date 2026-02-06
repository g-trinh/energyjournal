package firebase

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"energyjournal/internal/domain/user"

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
	apiKey string
}

func NewAuthProvider(client *Client, apiKey string) *AuthProvider {
	return &AuthProvider{client: client, apiKey: apiKey}
}

func (p *AuthProvider) CreateUser(ctx context.Context, email, password string) (string, error) {
	record, err := p.client.CreateUser(ctx, email, password)
	if err != nil {
		return "", err
	}
	return record.UID, nil
}

func (p *AuthProvider) Login(ctx context.Context, email, password string) (*user.AuthTokens, string, error) {
	endpoint := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", url.QueryEscape(p.apiKey))

	body, err := json.Marshal(map[string]any{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("marshal login request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, "", fmt.Errorf("create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, "", fmt.Errorf("firebase login failed: %s", errResp.Error.Message)
	}

	var result struct {
		IDToken      string `json:"idToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresIn    string `json:"expiresIn"`
		LocalID      string `json:"localId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", fmt.Errorf("decode login response: %w", err)
	}

	tokens := &user.AuthTokens{
		IDToken:      result.IDToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}
	return tokens, result.LocalID, nil
}

func (p *AuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*user.AuthTokens, string, error) {
	endpoint := fmt.Sprintf("https://securetoken.googleapis.com/v1/token?key=%s", url.QueryEscape(p.apiKey))

	body, err := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	})
	if err != nil {
		return nil, "", fmt.Errorf("marshal refresh request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, "", fmt.Errorf("create refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, "", fmt.Errorf("firebase refresh failed: %s", errResp.Error.Message)
	}

	var result struct {
		IDToken      string `json:"id_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    string `json:"expires_in"`
		UserID       string `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", fmt.Errorf("decode refresh response: %w", err)
	}

	tokens := &user.AuthTokens{
		IDToken:      result.IDToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}
	return tokens, result.UserID, nil
}
