package storage

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"

	"energyjournal/internal/domain/user"
	pkgerror "energyjournal/internal/pkg/error"
)

const activationTokensCollection = "activation_tokens"

type ActivationTokenRepository struct {
	client *firestore.Client
}

func NewActivationTokenRepository(client *firestore.Client) *ActivationTokenRepository {
	return &ActivationTokenRepository{client: client}
}

func (r *ActivationTokenRepository) Create(ctx context.Context, token *user.ActivationToken) error {
	_, err := r.client.Collection(activationTokensCollection).Doc(token.Token).Set(ctx, map[string]any{
		"token":     token.Token,
		"uid":       token.UID,
		"expiresAt": token.ExpiresAt,
	})
	return err
}

func (r *ActivationTokenRepository) GetByToken(ctx context.Context, token string) (*user.ActivationToken, error) {
	doc, err := r.client.Collection(activationTokensCollection).Doc(token).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, pkgerror.NewNotFoundError("activation_token", token)
		}
		return nil, err
	}

	return docToActivationToken(doc)
}

func (r *ActivationTokenRepository) Delete(ctx context.Context, token string) error {
	_, err := r.client.Collection(activationTokensCollection).Doc(token).Delete(ctx)
	return err
}

func (r *ActivationTokenRepository) FindExpired(ctx context.Context) ([]*user.ActivationToken, error) {
	now := time.Now()
	iter := r.client.Collection(activationTokensCollection).Where("expiresAt", "<", now).Documents(ctx)
	defer iter.Stop()

	var tokens []*user.ActivationToken
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		token, err := docToActivationToken(doc)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func docToActivationToken(doc *firestore.DocumentSnapshot) (*user.ActivationToken, error) {
	data := doc.Data()

	token := &user.ActivationToken{
		Token: getString(data, "token"),
		UID:   getString(data, "uid"),
	}

	if t, err := getTimestamp(data, "expiresAt"); err == nil {
		token.ExpiresAt = t
	}

	return token, nil
}
