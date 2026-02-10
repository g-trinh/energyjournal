package user

import (
	"encoding/json"
	"log"
	"net/http"

	"energyjournal/internal/domain/user"
	"energyjournal/internal/pkg/httputil"
	"energyjournal/internal/server/middleware"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create handles POST /users.
// Returns the same accepted response even when the email already exists
// to prevent account enumeration.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Invalid request body."})
		return
	}

	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Email, password, and confirmPassword are required."})
		return
	}

	if req.Password != req.ConfirmPassword {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Password and confirmPassword must match."})
		return
	}

	_, err := h.userService.Create(r.Context(), req.Email, req.Password, req.FirstName, req.LastName, req.Timezone)
	if err != nil {
		// Anti-enumeration: return the same success response regardless of error cause.
		// Log the actual error for operational visibility.
		log.Printf("user create error (suppressed for anti-enumeration): %v", err)
	}

	writeJSON(w, http.StatusCreated, CreateUserAcceptedResponse{
		Message: "Check your email to activate your account.",
		Status:  "pending_activation",
	})
}

// Activate handles POST /users/activate?token=...
func (h *UserHandler) Activate(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Missing activation token."})
		return
	}

	if err := h.userService.Activate(r.Context(), token); err != nil {
		statusCode, _ := httputil.MapErrors(err)
		writeJSON(w, statusCode, GenericErrorResponse{Message: "Activation failed."})
		return
	}

	writeJSON(w, http.StatusOK, ActivationResponse{Message: "Account activated successfully."})
}

// Login handles POST /users/login.
// All failure causes return the same generic message.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Invalid request body."})
		return
	}

	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Email and password are required."})
		return
	}

	tokens, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, GenericErrorResponse{Message: "Invalid email or password."})
		return
	}

	writeJSON(w, http.StatusOK, NewAuthTokensResponse(tokens))
}

// RefreshToken handles POST /users/refresh.
// Failures return a generic error response.
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Invalid request body."})
		return
	}

	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Refresh token is required."})
		return
	}

	tokens, err := h.userService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, GenericErrorResponse{Message: "Unable to refresh token."})
		return
	}

	writeJSON(w, http.StatusOK, NewAuthTokensResponse(tokens))
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, NewUserResponse(u))
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.userService.Update(r.Context(), u.UID, req.FirstName, req.LastName, req.Timezone)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, NewUserResponse(updated))
}

func (h *UserHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.userService.Delete(r.Context(), u.UID); err != nil {
		httputil.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
