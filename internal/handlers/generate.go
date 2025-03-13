package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
)

func generateResume(c *gin.Context) {
	var request models.GenerateResumeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// TODO: Implement resume generation logic
	response := models.GenerateResumeResponse{
		ResumeContent: models.ResumeContent{
			PersonalInfo: models.PersonalInfo{
				Name:  request.PersonalInfo.Name,
				Email: request.PersonalInfo.Email,
				Phone: request.PersonalInfo.Phone,
			},
			Summary:    "Generated professional summary based on job requirements",
			Experience: []models.Experience{},
			Education:  []models.Education{},
			Skills:     []string{},
			Projects:   []models.Project{},
		},
		Suggestions: []string{
			"Add more specific examples of your work",
			"Quantify your achievements",
		},
	}

	c.JSON(http.StatusOK, response)
}
