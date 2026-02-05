package user

import (
	"encoding/json"
	"energyjournal/internal/server/middleware"
	"net/http"

	"energyjournal/internal/domain/user"
	"energyjournal/internal/pkg/httputil"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	u, err := h.userService.Create(r.Context(), req.Email, req.Password, req.FirstName, req.LastName, req.Timezone)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewUserResponse(u))
}

func (h *UserHandler) Activate(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	if err := h.userService.Activate(r.Context(), token); err != nil {
		httputil.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewUserResponse(u))
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewUserResponse(updated))
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
