package user

import "context"

type UserService interface {
	Create(ctx context.Context, email, password, firstname, lastname, timezone string) (*User, error)
	Activate(ctx context.Context, token string) error
	GetByUID(ctx context.Context, uid string) (*User, error)
	Update(ctx context.Context, uid, firstname, lastname, timezone string) (*User, error)
	Delete(ctx context.Context, uid string) error
	CleanupExpired(ctx context.Context) error
	Login(ctx context.Context, email, password string) (*AuthTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error)
}

type AuthProvider interface {
	CreateUser(ctx context.Context, email, password string) (uid string, err error)
	Login(ctx context.Context, email, password string) (*AuthTokens, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, string, error)
}

type EmailSender interface {
	SendActivationEmail(ctx context.Context, email, activationLink string) error
}
