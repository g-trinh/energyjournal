package user

import (
	"time"

	"energyjournal/internal/domain/user"
)

type UserResponse struct {
	UID       string    `json:"uid"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Timezone  string    `json:"timezone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type AuthTokensResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

func NewAuthTokensResponse(t *user.AuthTokens) *AuthTokensResponse {
	return &AuthTokensResponse{
		IDToken:      t.IDToken,
		RefreshToken: t.RefreshToken,
		ExpiresIn:    t.ExpiresIn,
	}
}

func NewUserResponse(u *user.User) *UserResponse {
	return &UserResponse{
		UID:       u.UID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Timezone:  u.Timezone,
		Status:    string(u.Status),
		CreatedAt: u.CreatedAt,
	}
}
