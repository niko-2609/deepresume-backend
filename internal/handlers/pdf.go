package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
	"github.com/nikolai/ai-resume-builder/backend/internal/utils"
)

func generatePDF(c *gin.Context) {
	var request models.ResumeContent

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// TODO: Implement PDF generation logic
	pdfBytes, err := utils.GeneratePDF(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate PDF",
		})
		return
	}

	// Set response headers for PDF download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=resume.pdf")
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
