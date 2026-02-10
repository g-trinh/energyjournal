package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"energyjournal/internal/domain/user"
	pkgerror "energyjournal/internal/pkg/error"
)

// --- Mock service ---

type mockUserService struct {
	createFn       func(ctx context.Context, email, password, firstname, lastname, timezone string) (*user.User, error)
	activateFn     func(ctx context.Context, token string) error
	loginFn        func(ctx context.Context, email, password string) (*user.AuthTokens, error)
	refreshTokenFn func(ctx context.Context, refreshToken string) (*user.AuthTokens, error)
}

func (m *mockUserService) Create(ctx context.Context, email, password, firstname, lastname, timezone string) (*user.User, error) {
	return m.createFn(ctx, email, password, firstname, lastname, timezone)
}

func (m *mockUserService) Activate(ctx context.Context, token string) error {
	return m.activateFn(ctx, token)
}

func (m *mockUserService) Login(ctx context.Context, email, password string) (*user.AuthTokens, error) {
	return m.loginFn(ctx, email, password)
}

func (m *mockUserService) RefreshToken(ctx context.Context, refreshToken string) (*user.AuthTokens, error) {
	return m.refreshTokenFn(ctx, refreshToken)
}

func (m *mockUserService) GetByUID(ctx context.Context, uid string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserService) Update(ctx context.Context, uid, firstname, lastname, timezone string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserService) Delete(ctx context.Context, uid string) error {
	return nil
}

func (m *mockUserService) CleanupExpired(ctx context.Context) error {
	return nil
}

// --- Helper ---

func postJSON(t *testing.T, handler http.HandlerFunc, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

func decodeJSON[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(rr.Body).Decode(&v); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return v
}

// --- Tests ---

// Anti-enumeration: signup with new email returns accepted response.
func TestCreate_Success_ReturnsAcceptedResponse(t *testing.T) {
	svc := &mockUserService{
		createFn: func(ctx context.Context, email, password, firstname, lastname, timezone string) (*user.User, error) {
			return &user.User{
				UID:       "uid-1",
				Email:     email,
				Status:    user.StatusPendingValidation,
				CreatedAt: time.Now(),
			}, nil
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Create, "/users", map[string]string{
		"email":           "new@example.com",
		"password":        "secret123",
		"confirmPassword": "secret123",
	})

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}

	resp := decodeJSON[CreateUserAcceptedResponse](t, rr)
	if resp.Message != "Check your email to activate your account." {
		t.Errorf("unexpected message: %s", resp.Message)
	}
	if resp.Status != "pending_activation" {
		t.Errorf("unexpected status: %s", resp.Status)
	}
}

// Anti-enumeration: signup with existing email returns the same accepted response.
func TestCreate_DuplicateEmail_ReturnsSameAcceptedResponse(t *testing.T) {
	svc := &mockUserService{
		createFn: func(ctx context.Context, email, password, firstname, lastname, timezone string) (*user.User, error) {
			return nil, errors.New("email already exists")
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Create, "/users", map[string]string{
		"email":           "existing@example.com",
		"password":        "secret123",
		"confirmPassword": "secret123",
	})

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 even for duplicate email, got %d", rr.Code)
	}

	resp := decodeJSON[CreateUserAcceptedResponse](t, rr)
	if resp.Message != "Check your email to activate your account." {
		t.Errorf("unexpected message: %s", resp.Message)
	}
	if resp.Status != "pending_activation" {
		t.Errorf("unexpected status: %s", resp.Status)
	}
}

// Anti-enumeration: any internal error during signup still returns accepted response.
func TestCreate_InternalError_ReturnsSameAcceptedResponse(t *testing.T) {
	svc := &mockUserService{
		createFn: func(ctx context.Context, email, password, firstname, lastname, timezone string) (*user.User, error) {
			return nil, errors.New("firestore write failed")
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Create, "/users", map[string]string{
		"email":           "test@example.com",
		"password":        "secret123",
		"confirmPassword": "secret123",
	})

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 even for internal error, got %d", rr.Code)
	}

	resp := decodeJSON[CreateUserAcceptedResponse](t, rr)
	if resp.Message != "Check your email to activate your account." {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

// Create: missing required fields returns 400.
func TestCreate_MissingFields_Returns400(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Create, "/users", map[string]string{
		"email":    "test@example.com",
		"password": "secret123",
		// missing confirmPassword
	})

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message == "" {
		t.Error("expected error message in response")
	}
}

// Create: password mismatch returns 400.
func TestCreate_PasswordMismatch_Returns400(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Create, "/users", map[string]string{
		"email":           "test@example.com",
		"password":        "secret123",
		"confirmPassword": "different456",
	})

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Password and confirmPassword must match." {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

// Login: unknown email returns generic 401.
func TestLogin_UnknownEmail_ReturnsGeneric401(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(ctx context.Context, email, password string) (*user.AuthTokens, error) {
			return nil, errors.New("firebase login failed: EMAIL_NOT_FOUND")
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Login, "/users/login", map[string]string{
		"email":    "unknown@example.com",
		"password": "secret123",
	})

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Invalid email or password." {
		t.Errorf("expected generic message, got: %s", resp.Message)
	}
}

// Login: wrong password returns generic 401.
func TestLogin_WrongPassword_ReturnsGeneric401(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(ctx context.Context, email, password string) (*user.AuthTokens, error) {
			return nil, errors.New("firebase login failed: INVALID_PASSWORD")
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Login, "/users/login", map[string]string{
		"email":    "valid@example.com",
		"password": "wrong",
	})

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Invalid email or password." {
		t.Errorf("expected generic message, got: %s", resp.Message)
	}
}

// Login: inactive account returns generic 401 (same as wrong password).
func TestLogin_InactiveAccount_ReturnsGeneric401(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(ctx context.Context, email, password string) (*user.AuthTokens, error) {
			return nil, pkgerror.NewInputValidationError("user", "user is not active")
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Login, "/users/login", map[string]string{
		"email":    "inactive@example.com",
		"password": "correct",
	})

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Invalid email or password." {
		t.Errorf("expected generic message, got: %s", resp.Message)
	}
}

// Login: deleted account returns generic 401.
func TestLogin_DeletedAccount_ReturnsGeneric401(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(ctx context.Context, email, password string) (*user.AuthTokens, error) {
			return nil, pkgerror.NewInputValidationError("user", "user is not active")
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Login, "/users/login", map[string]string{
		"email":    "deleted@example.com",
		"password": "anything",
	})

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Invalid email or password." {
		t.Errorf("expected generic message, got: %s", resp.Message)
	}
}

// Login: success returns 200 with tokens.
func TestLogin_Success_ReturnsTokens(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(ctx context.Context, email, password string) (*user.AuthTokens, error) {
			return &user.AuthTokens{
				IDToken:      "id-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    "3600",
			}, nil
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Login, "/users/login", map[string]string{
		"email":    "active@example.com",
		"password": "correct",
	})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	resp := decodeJSON[AuthTokensResponse](t, rr)
	if resp.IDToken != "id-token" {
		t.Errorf("unexpected idToken: %s", resp.IDToken)
	}
	if resp.RefreshToken != "refresh-token" {
		t.Errorf("unexpected refreshToken: %s", resp.RefreshToken)
	}
}

// Login: missing fields returns 400.
func TestLogin_MissingFields_Returns400(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.Login, "/users/login", map[string]string{
		"email": "test@example.com",
	})

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

// Activate: success returns 200 with activation response.
func TestActivate_Success_ReturnsActivationResponse(t *testing.T) {
	svc := &mockUserService{
		activateFn: func(ctx context.Context, token string) error {
			return nil
		},
	}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/users/activate?token=valid-token", nil)
	rr := httptest.NewRecorder()
	h.Activate(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	resp := decodeJSON[ActivationResponse](t, rr)
	if resp.Message != "Account activated successfully." {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

// Activate: missing token returns 400 with generic error.
func TestActivate_MissingToken_Returns400(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/users/activate", nil)
	rr := httptest.NewRecorder()
	h.Activate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Missing activation token." {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

// Activate: expired token returns generic activation error.
func TestActivate_ExpiredToken_ReturnsGenericError(t *testing.T) {
	svc := &mockUserService{
		activateFn: func(ctx context.Context, token string) error {
			return pkgerror.NewInputValidationError("token", "token has expired")
		},
	}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/users/activate?token=expired-token", nil)
	rr := httptest.NewRecorder()
	h.Activate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Activation failed." {
		t.Errorf("expected generic activation error, got: %s", resp.Message)
	}
}

// Activate: unknown token returns generic activation error.
func TestActivate_UnknownToken_ReturnsGenericError(t *testing.T) {
	svc := &mockUserService{
		activateFn: func(ctx context.Context, token string) error {
			return pkgerror.NewNotFoundError("activation_token", token)
		},
	}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/users/activate?token=nonexistent", nil)
	rr := httptest.NewRecorder()
	h.Activate(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Activation failed." {
		t.Errorf("expected generic activation error, got: %s", resp.Message)
	}
}

// Activate: already-active user returns generic activation error.
func TestActivate_AlreadyActive_ReturnsGenericError(t *testing.T) {
	svc := &mockUserService{
		activateFn: func(ctx context.Context, token string) error {
			return pkgerror.NewInputValidationError("user", "user is already active")
		},
	}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/users/activate?token=used-token", nil)
	rr := httptest.NewRecorder()
	h.Activate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Activation failed." {
		t.Errorf("expected generic activation error, got: %s", resp.Message)
	}
}

// Refresh: success returns 200 with new tokens.
func TestRefreshToken_Success_ReturnsTokens(t *testing.T) {
	svc := &mockUserService{
		refreshTokenFn: func(ctx context.Context, refreshToken string) (*user.AuthTokens, error) {
			return &user.AuthTokens{
				IDToken:      "new-id-token",
				RefreshToken: "new-refresh-token",
				ExpiresIn:    "3600",
			}, nil
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.RefreshToken, "/users/refresh", map[string]string{
		"refreshToken": "valid-refresh-token",
	})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	resp := decodeJSON[AuthTokensResponse](t, rr)
	if resp.IDToken != "new-id-token" {
		t.Errorf("unexpected idToken: %s", resp.IDToken)
	}
}

// Refresh: invalid token returns generic error.
func TestRefreshToken_InvalidToken_ReturnsGenericError(t *testing.T) {
	svc := &mockUserService{
		refreshTokenFn: func(ctx context.Context, refreshToken string) (*user.AuthTokens, error) {
			return nil, errors.New("firebase refresh failed: TOKEN_EXPIRED")
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.RefreshToken, "/users/refresh", map[string]string{
		"refreshToken": "invalid-token",
	})

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Unable to refresh token." {
		t.Errorf("expected generic message, got: %s", resp.Message)
	}
}

// Refresh: missing token returns 400.
func TestRefreshToken_MissingToken_Returns400(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.RefreshToken, "/users/refresh", map[string]string{})

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Refresh token is required." {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

// Refresh: inactive user returns generic error.
func TestRefreshToken_InactiveUser_ReturnsGenericError(t *testing.T) {
	svc := &mockUserService{
		refreshTokenFn: func(ctx context.Context, refreshToken string) (*user.AuthTokens, error) {
			return nil, pkgerror.NewInputValidationError("user", "user is not active")
		},
	}
	h := NewUserHandler(svc)

	rr := postJSON(t, h.RefreshToken, "/users/refresh", map[string]string{
		"refreshToken": "valid-but-inactive-user",
	})

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	resp := decodeJSON[GenericErrorResponse](t, rr)
	if resp.Message != "Unable to refresh token." {
		t.Errorf("expected generic message, got: %s", resp.Message)
	}
}

// Verify all auth error responses use JSON Content-Type.
func TestAllErrorResponses_UseJSONContentType(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(ctx context.Context, email, password string) (*user.AuthTokens, error) {
			return nil, errors.New("any error")
		},
		refreshTokenFn: func(ctx context.Context, refreshToken string) (*user.AuthTokens, error) {
			return nil, errors.New("any error")
		},
		activateFn: func(ctx context.Context, token string) error {
			return errors.New("any error")
		},
	}
	h := NewUserHandler(svc)

	// Login error
	rr := postJSON(t, h.Login, "/users/login", map[string]string{
		"email":    "test@example.com",
		"password": "pass",
	})
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("login error Content-Type: %s", ct)
	}

	// Refresh error
	rr = postJSON(t, h.RefreshToken, "/users/refresh", map[string]string{
		"refreshToken": "bad",
	})
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("refresh error Content-Type: %s", ct)
	}

	// Activate error
	req := httptest.NewRequest(http.MethodPost, "/users/activate?token=bad", nil)
	rr = httptest.NewRecorder()
	h.Activate(rr, req)
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("activate error Content-Type: %s", ct)
	}
}
