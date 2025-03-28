package handlers

import (
	"github.com/nikolai/ai-resume-builder/backend/internal/service"
)

// Handler contains all dependencies needed for the HTTP handlers
type BaseHandler struct {
	service *service.BaseService
}

// NewHandler creates a new Handler instance
func NewBaseHandler(service *service.BaseService) *BaseHandler {
	return &BaseHandler{
		service: service,
	}
}
