package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
	"github.com/nikolai/ai-resume-builder/backend/internal/service"
)

type OnboardingRequest struct {
	User           models.User             `json:"user" binding:"required"`
	WorkExperience []models.WorkExperience `json:"workExperience" binding:"required"`
	Education      []models.Education      `json:"education" binding:"required"`
}

// HandleOnboarding godoc
// @Summary Complete user onboarding
// @Description Process complete onboarding data including user profile, work experience, and education
// @Tags onboarding
// @Accept json
// @Produce json
// @Param request body OnboardingRequest true "Onboarding Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/onboarding [post]
func (h *UserHandler) HandleOnboarding(c *gin.Context) {
	var req OnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Start a transaction to save all the data
	err := h.userService.CreateUserOnboarding(c.Request.Context(), service.OnboardingData{
		User:           req.User,
		WorkExperience: req.WorkExperience,
		Education:      req.Education,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save onboarding data: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Onboarding completed successfully",
	})
}
