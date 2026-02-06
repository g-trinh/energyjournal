package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByUID(ctx context.Context, uid string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type ActivationTokenRepository interface {
	Create(ctx context.Context, token *ActivationToken) error
	GetByToken(ctx context.Context, token string) (*ActivationToken, error)
	Delete(ctx context.Context, token string) error
	FindExpired(ctx context.Context) ([]*ActivationToken, error)
}
