package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikolai/ai-resume-builder/backend/internal/service"
)

type ResumeHandler struct {
	resumeService *service.ResumeService
}

type GenerateResumeRequest struct {
	UserID         uint   `json:"userId" binding:"required"`
	JobDescription string `json:"jobDescription" binding:"required"`
}

type StreamChunk struct {
	Chunk string `json:"chunk"`
	Done  bool   `json:"done"`
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

	// Set headers for streaming
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Create a flusher to ensure data is sent immediately
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// Stream the resume generation
	err := h.resumeService.GenerateResume(c.Request.Context(), req.UserID, req.JobDescription, func(chunk string, done bool) error {
		// Create the chunk response
		chunkResponse := StreamChunk{
			Chunk: chunk,
			Done:  done,
		}

		// Convert to JSON
		data, err := json.Marshal(chunkResponse)
		if err != nil {
			return err
		}

		// Send the chunk to the client
		_, err = c.Writer.Write(data)
		if err != nil {
			return err
		}
		c.Writer.Write([]byte("\n"))
		flusher.Flush()

		// If we're done, send the source information
		if done {
			// Send source information as the final message
			sourceData, _ := json.Marshal(map[string]string{"source": "llm"})
			c.Writer.Write(sourceData)
			c.Writer.Write([]byte("\n"))
			flusher.Flush()
		}

		return nil
	})

	if err != nil {
		// If there's an error after we've started streaming, we can't use regular error responses
		// Just log the error and close the connection
		errorData, _ := json.Marshal(map[string]string{"error": err.Error()})
		c.Writer.Write(errorData)
		c.Writer.Write([]byte("\n"))
		flusher.Flush()
	}
}
