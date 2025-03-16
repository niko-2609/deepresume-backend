package handlers

import (
	"github.com/nikolai/ai-resume-builder/backend/internal/service"
)

// Handler contains all dependencies needed for the HTTP handlers
type Handler struct {
	service *service.Service
}

// NewHandler creates a new Handler instance
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
