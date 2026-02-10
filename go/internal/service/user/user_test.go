package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"energyjournal/internal/domain/user"
)

// --- Mock repositories and providers ---

type mockUserRepo struct {
	users map[string]*user.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*user.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	m.users[u.UID] = u
	return nil
}

func (m *mockUserRepo) GetByUID(ctx context.Context, uid string) (*user.User, error) {
	u, ok := m.users[uid]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	m.users[u.UID] = u
	return nil
}

type mockTokenRepo struct {
	tokens map[string]*user.ActivationToken
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[string]*user.ActivationToken)}
}

func (m *mockTokenRepo) Create(ctx context.Context, token *user.ActivationToken) error {
	m.tokens[token.Token] = token
	return nil
}

func (m *mockTokenRepo) GetByToken(ctx context.Context, token string) (*user.ActivationToken, error) {
	t, ok := m.tokens[token]
	if !ok {
		return nil, errors.New("token not found")
	}
	return t, nil
}

func (m *mockTokenRepo) Delete(ctx context.Context, token string) error {
	delete(m.tokens, token)
	return nil
}

func (m *mockTokenRepo) FindExpired(ctx context.Context) ([]*user.ActivationToken, error) {
	var expired []*user.ActivationToken
	for _, t := range m.tokens {
		if t.IsExpired() {
			expired = append(expired, t)
		}
	}
	return expired, nil
}

type mockAuthProvider struct {
	createErr error
	loginErr  error
}

func (m *mockAuthProvider) CreateUser(ctx context.Context, email, password string) (string, error) {
	if m.createErr != nil {
		return "", m.createErr
	}
	return "uid-123", nil
}

func (m *mockAuthProvider) Login(ctx context.Context, email, password string) (*user.AuthTokens, string, error) {
	if m.loginErr != nil {
		return nil, "", m.loginErr
	}
	return &user.AuthTokens{
		IDToken:      "id-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    "3600",
	}, "uid-123", nil
}

func (m *mockAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*user.AuthTokens, string, error) {
	return nil, "", errors.New("not implemented in test")
}

type mockEmailSender struct {
	lastLink string
}

func (m *mockEmailSender) SendActivationEmail(ctx context.Context, email, activationLink string) error {
	m.lastLink = activationLink
	return nil
}

// --- Tests ---

// Activation: expired token is rejected.
func TestActivate_ExpiredToken_Rejected(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()

	userRepo.users["uid-1"] = &user.User{
		UID:    "uid-1",
		Email:  "test@example.com",
		Status: user.StatusPendingValidation,
	}

	tokenRepo.tokens["expired-token"] = &user.ActivationToken{
		Token:     "expired-token",
		UID:       "uid-1",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // expired
	}

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, &mockEmailSender{}, "http://localhost:8080")

	err := svc.Activate(context.Background(), "expired-token")
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

// Activation: valid token activates the user.
func TestActivate_ValidToken_ActivatesUser(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()

	userRepo.users["uid-1"] = &user.User{
		UID:    "uid-1",
		Email:  "test@example.com",
		Status: user.StatusPendingValidation,
	}

	tokenRepo.tokens["valid-token"] = &user.ActivationToken{
		Token:     "valid-token",
		UID:       "uid-1",
		ExpiresAt: time.Now().Add(23 * time.Hour), // within 24h window
	}

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, &mockEmailSender{}, "http://localhost:8080")

	err := svc.Activate(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u := userRepo.users["uid-1"]
	if u.Status != user.StatusActive {
		t.Errorf("expected user status ACTIVE, got %s", u.Status)
	}
}

// Activation: token is consumed (one-time use).
func TestActivate_TokenConsumed_AfterUse(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()

	userRepo.users["uid-1"] = &user.User{
		UID:    "uid-1",
		Email:  "test@example.com",
		Status: user.StatusPendingValidation,
	}

	tokenRepo.tokens["one-time-token"] = &user.ActivationToken{
		Token:     "one-time-token",
		UID:       "uid-1",
		ExpiresAt: time.Now().Add(23 * time.Hour),
	}

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, &mockEmailSender{}, "http://localhost:8080")

	// First use: succeeds
	err := svc.Activate(context.Background(), "one-time-token")
	if err != nil {
		t.Fatalf("first activation failed: %v", err)
	}

	// Second use: fails (token deleted)
	err = svc.Activate(context.Background(), "one-time-token")
	if err == nil {
		t.Fatal("expected error on second activation (token consumed), got nil")
	}
}

// Activation: unknown token is rejected.
func TestActivate_UnknownToken_Rejected(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, &mockEmailSender{}, "http://localhost:8080")

	err := svc.Activate(context.Background(), "nonexistent-token")
	if err == nil {
		t.Fatal("expected error for unknown token, got nil")
	}
}

// Login: inactive user is rejected.
func TestLogin_InactiveUser_Rejected(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()

	userRepo.users["uid-123"] = &user.User{
		UID:    "uid-123",
		Email:  "pending@example.com",
		Status: user.StatusPendingValidation,
	}

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, &mockEmailSender{}, "http://localhost:8080")

	_, err := svc.Login(context.Background(), "pending@example.com", "password")
	if err == nil {
		t.Fatal("expected error for inactive user, got nil")
	}
}

// Login: active user succeeds.
func TestLogin_ActiveUser_Succeeds(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()

	userRepo.users["uid-123"] = &user.User{
		UID:    "uid-123",
		Email:  "active@example.com",
		Status: user.StatusActive,
	}

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, &mockEmailSender{}, "http://localhost:8080")

	tokens, err := svc.Login(context.Background(), "active@example.com", "password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens.IDToken == "" {
		t.Error("expected non-empty IDToken")
	}
}

// Login: deleted user is rejected.
func TestLogin_DeletedUser_Rejected(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()

	userRepo.users["uid-123"] = &user.User{
		UID:    "uid-123",
		Email:  "deleted@example.com",
		Status: user.StatusDeleted,
	}

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, &mockEmailSender{}, "http://localhost:8080")

	_, err := svc.Login(context.Background(), "deleted@example.com", "password")
	if err == nil {
		t.Fatal("expected error for deleted user, got nil")
	}
}

// Create: activation link contains the configured base URL and token.
func TestCreate_ActivationLink_Format(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()
	emailSender := &mockEmailSender{}

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, emailSender, "https://app.example.com")

	_, err := svc.Create(context.Background(), "test@example.com", "password123", "", "", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if emailSender.lastLink == "" {
		t.Fatal("expected activation link to be sent")
	}

	// Link should start with the configured base URL
	expectedPrefix := "https://app.example.com/activate?token="
	if len(emailSender.lastLink) <= len(expectedPrefix) {
		t.Fatalf("activation link too short: %s", emailSender.lastLink)
	}
	if emailSender.lastLink[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected link to start with %s, got: %s", expectedPrefix, emailSender.lastLink)
	}
}

// Create: activation token expires after 24 hours.
func TestCreate_ActivationToken_24HourExpiry(t *testing.T) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()

	svc := NewUserService(userRepo, tokenRepo, &mockAuthProvider{}, &mockEmailSender{}, "http://localhost:8080")

	_, err := svc.Create(context.Background(), "test@example.com", "password123", "", "", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find the created token
	if len(tokenRepo.tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokenRepo.tokens))
	}

	for _, tok := range tokenRepo.tokens {
		expiresIn := time.Until(tok.ExpiresAt)
		if expiresIn < 23*time.Hour || expiresIn > 25*time.Hour {
			t.Errorf("expected token to expire in ~24 hours, expires in %v", expiresIn)
		}
	}
}
