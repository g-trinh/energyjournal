package user

import "time"

type UserStatus string

const (
	StatusPendingValidation UserStatus = "PENDING_VALIDATION"
	StatusActive            UserStatus = "ACTIVE"
	StatusDeleted           UserStatus = "DELETED"
)

type User struct {
	UID       string
	Email     string
	FirstName string
	LastName  string
	Timezone  string
	Status    UserStatus
	CreatedAt time.Time
	DeletedAt *time.Time
}

type ActivationToken struct {
	Token     string
	UID       string
	ExpiresAt time.Time
}

func (t *ActivationToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
