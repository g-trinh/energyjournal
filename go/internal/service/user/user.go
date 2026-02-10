package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"energyjournal/internal/domain/user"
	pkgerror "energyjournal/internal/pkg/error"
)

type userService struct {
	userRepo              user.UserRepository
	tokenRepo             user.ActivationTokenRepository
	authProvider          user.AuthProvider
	emailSender           user.EmailSender
	activationBaseURL     string
}

func NewUserService(userRepo user.UserRepository, tokenRepo user.ActivationTokenRepository, authProvider user.AuthProvider, emailSender user.EmailSender, activationBaseURL string) user.UserService {
	return &userService{
		userRepo:          userRepo,
		tokenRepo:         tokenRepo,
		authProvider:      authProvider,
		emailSender:       emailSender,
		activationBaseURL: activationBaseURL,
	}
}

func (s *userService) Create(ctx context.Context, email, password, firstname, lastname, timezone string) (*user.User, error) {
	// Create Firebase Auth user first
	uid, err := s.authProvider.CreateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		UID:       uid,
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Timezone:  timezone,
		Status:    user.StatusPendingValidation,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	activationToken := &user.ActivationToken{
		Token:     token,
		UID:       uid,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.tokenRepo.Create(ctx, activationToken); err != nil {
		return nil, err
	}

	activationLink := fmt.Sprintf("%s/activate?token=%s", s.activationBaseURL, token)
	if err := s.emailSender.SendActivationEmail(ctx, email, activationLink); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) Activate(ctx context.Context, token string) error {
	activationToken, err := s.tokenRepo.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	if activationToken.IsExpired() {
		return pkgerror.NewInputValidationError("token", "token has expired")
	}

	u, err := s.userRepo.GetByUID(ctx, activationToken.UID)
	if err != nil {
		return err
	}

	if u.Status == user.StatusActive {
		return pkgerror.NewInputValidationError("user", "user is already active")
	}

	u.Status = user.StatusActive
	if err := s.userRepo.Update(ctx, u); err != nil {
		return err
	}

	return s.tokenRepo.Delete(ctx, token)
}

func (s *userService) GetByUID(ctx context.Context, uid string) (*user.User, error) {
	return s.userRepo.GetByUID(ctx, uid)
}

func (s *userService) Update(ctx context.Context, uid, firstname, lastname, timezone string) (*user.User, error) {
	u, err := s.userRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	u.FirstName = firstname
	u.LastName = lastname
	u.Timezone = timezone

	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) Delete(ctx context.Context, uid string) error {
	u, err := s.userRepo.GetByUID(ctx, uid)
	if err != nil {
		return err
	}

	now := time.Now()
	u.Status = user.StatusDeleted
	u.DeletedAt = &now

	return s.userRepo.Update(ctx, u)
}

func (s *userService) CleanupExpired(ctx context.Context) error {
	expiredTokens, err := s.tokenRepo.FindExpired(ctx)
	if err != nil {
		return err
	}

	for _, token := range expiredTokens {
		u, err := s.userRepo.GetByUID(ctx, token.UID)
		if err != nil {
			continue
		}

		if u.Status == user.StatusPendingValidation {
			now := time.Now()
			u.Status = user.StatusDeleted
			u.DeletedAt = &now
			_ = s.userRepo.Update(ctx, u)
		}

		_ = s.tokenRepo.Delete(ctx, token.Token)
	}

	return nil
}

func (s *userService) Login(ctx context.Context, email, password string) (*user.AuthTokens, error) {
	tokens, uid, err := s.authProvider.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if u.Status != user.StatusActive {
		return nil, pkgerror.NewInputValidationError("user", "user is not active")
	}

	return tokens, nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*user.AuthTokens, error) {
	tokens, uid, err := s.authProvider.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if u.Status != user.StatusActive {
		return nil, pkgerror.NewInputValidationError("user", "user is not active")
	}

	return tokens, nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
