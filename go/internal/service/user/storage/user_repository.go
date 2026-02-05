package storage

import (
	"context"

	"cloud.google.com/go/firestore"

	"energyjournal/internal/domain/user"
	pkgerror "energyjournal/internal/pkg/error"
)

const usersCollection = "users"

type UserRepository struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	_, err := r.client.Collection(usersCollection).Doc(u.UID).Set(ctx, map[string]interface{}{
		"uid":       u.UID,
		"email":     u.Email,
		"firstname": u.FirstName,
		"lastname":  u.LastName,
		"timezone":  u.Timezone,
		"status":    string(u.Status),
		"createdAt": u.CreatedAt,
		"deletedAt": u.DeletedAt,
	})
	return err
}

func (r *UserRepository) GetByUID(ctx context.Context, uid string) (*user.User, error) {
	doc, err := r.client.Collection(usersCollection).Doc(uid).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, pkgerror.NewNotFoundError("user", uid)
		}
		return nil, err
	}

	return docToUser(doc)
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	_, err := r.client.Collection(usersCollection).Doc(u.UID).Set(ctx, map[string]interface{}{
		"uid":       u.UID,
		"email":     u.Email,
		"firstname": u.FirstName,
		"lastname":  u.LastName,
		"timezone":  u.Timezone,
		"status":    string(u.Status),
		"createdAt": u.CreatedAt,
		"deletedAt": u.DeletedAt,
	})
	return err
}

func docToUser(doc *firestore.DocumentSnapshot) (*user.User, error) {
	data := doc.Data()

	u := &user.User{
		UID:       getString(data, "uid"),
		Email:     getString(data, "email"),
		FirstName: getString(data, "firstname"),
		LastName:  getString(data, "lastname"),
		Timezone:  getString(data, "timezone"),
		Status:    user.UserStatus(getString(data, "status")),
	}

	if t, err := getTimestamp(data, "createdAt"); err == nil {
		u.CreatedAt = t
	}

	if t, err := getTimestamp(data, "deletedAt"); err == nil && !t.IsZero() {
		u.DeletedAt = &t
	}

	return u, nil
}
