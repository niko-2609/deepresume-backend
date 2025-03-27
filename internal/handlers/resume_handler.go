package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikolai/ai-resume-builder/backend/internal/service"
)

type ResumeHandler struct {
	resumeService *service.ResumeService
}

type GenerateResumeRequest struct {
	JobDescription string `json:"jobDescription" binding:"required"`
}

func NewResumeHandler(resumeService *service.ResumeService) *ResumeHandler {
	return &ResumeHandler{resumeService: resumeService}
}

// GenerateResume handles requests to generate a resume from a job description
func (h *ResumeHandler) GenerateResume(c *gin.Context) {
	var req GenerateResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to generate the resume
	resume, err := h.resumeService.GenerateResume(c.Request.Context(), req.JobDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resume": resume,
	})
}
