package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
	"github.com/nikolai/ai-resume-builder/backend/internal/service"
)

// UserHandler handles all user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

type OnboardingRequest struct {
	User           models.User             `json:"user" binding:"required"`
	WorkExperience []models.WorkExperience `json:"workExperience" binding:"required"`
	Education      []models.Education      `json:"education" binding:"required"`
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.userService.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := h.userService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user.ID = uint(id)
	if err := h.userService.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
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
