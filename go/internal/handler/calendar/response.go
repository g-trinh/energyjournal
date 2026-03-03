package calendar

import "energyjournal/internal/domain/calendar"

type StatusResponse struct {
	Status calendar.ConnectionStatus `json:"status"`
}

type AuthURLResponse struct {
	AuthURL string `json:"auth_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
