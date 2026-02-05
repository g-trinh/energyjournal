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
