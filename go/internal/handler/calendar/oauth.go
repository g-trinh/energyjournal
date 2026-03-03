package calendar

import (
	"net/http"
	"net/url"

	"energyjournal/internal/domain/calendar"
	"energyjournal/internal/server/middleware"
)

type OAuthHandler struct {
	service         calendar.CalendarService
	frontendBaseURL string
}

func NewOAuthHandler(service calendar.CalendarService, frontendBaseURL string) *OAuthHandler {
	return &OAuthHandler{
		service:         service,
		frontendBaseURL: frontendBaseURL,
	}
}

// GetAuthURL godoc
// @Summary Build Google OAuth authorization URL
// @Tags calendar
// @Security BearerAuth
// @Success 200 {object} calendar.AuthURLResponse
// @Failure 401 {object} calendar.ErrorResponse
// @Failure 500 {object} calendar.ErrorResponse
// @Router /calendar/auth [get]
func (h *OAuthHandler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UIDFromContext(r.Context())
	if !ok || uid == "" {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	writeJSON(w, http.StatusOK, AuthURLResponse{AuthURL: h.service.BuildAuthURL(uid)})
}

// Callback godoc
// @Summary Google OAuth callback
// @Description Receives the authorization code from Google, validates state and stores tokens.
// @Tags calendar
// @Param code query string true "Authorization code"
// @Param state query string true "Signed state"
// @Success 302 {string} string "Redirect to frontend /timespending?calendar=select"
// @Failure 400 {object} calendar.ErrorResponse
// @Failure 500 {object} calendar.ErrorResponse
// @Router /calendar/auth/callback [get]
func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing code or state"})
		return
	}

	if err := h.service.HandleCallback(r.Context(), code, state); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "oauth callback failed"})
		return
	}

	redirectURL, err := url.JoinPath(h.frontendBaseURL, "timespending")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "invalid frontend base url"})
		return
	}
	http.Redirect(w, r, redirectURL+"?calendar=select", http.StatusFound)
}
