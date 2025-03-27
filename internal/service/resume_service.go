package service

import (
	"context"
	"fmt"

	"github.com/nikolai/ai-resume-builder/backend/internal/interfaces"
)

type ResumeService struct {
	db             interfaces.DB
	keywordService *KeywordService
}

func NewResumeService(db interfaces.DB, keywordService *KeywordService) *ResumeService {
	return &ResumeService{
		db:             db,
		keywordService: keywordService,
	}
}

// GenerateResume generates a resume based on the job description
func (s *ResumeService) GenerateResume(ctx context.Context, jobDescription string) (string, error) {
	// Extract keywords from job description
	keywords := s.keywordService.ExtractAndRankKeywords(jobDescription)

	// For the MVP, we'll just generate a simple Markdown-formatted resume
	// based on the extracted keywords
	keywordStrings := make([]string, 0, len(keywords))

	// Use up to 10 keywords, but handle cases where fewer are available
	maxKeywords := 10
	if len(keywords) < maxKeywords {
		maxKeywords = len(keywords)
	}

	// Safely get the top keywords
	for _, k := range keywords[:maxKeywords] {
		keywordStrings = append(keywordStrings, k.Word)
	}

	// Generate a simple resume with the extracted keywords
	// Ensure we don't try to access beyond the available keywords
	summaryKeywords := keywordStrings
	if len(summaryKeywords) > 3 {
		summaryKeywords = summaryKeywords[:3]
	}

	resume := fmt.Sprintf(`# Professional Resume

## Summary
Experienced professional with expertise in %s.

## Skills
`, joinWithCommas(summaryKeywords))

	// Add skills section
	for _, keyword := range keywordStrings {
		resume += fmt.Sprintf("- %s\n", keyword)
	}

	resume += `
## Experience
- Implemented solutions utilizing the latest technologies
- Developed and maintained web applications
- Collaborated with cross-functional teams

## Education
- Bachelor's Degree in Computer Science

*Note: This is an automatically generated resume based on the job description. Please customize it with your actual experience and credentials.*
`

	return resume, nil
}

// Helper function to join strings with commas and "and" for the last item
func joinWithCommas(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 {
		return items[0] + " and " + items[1]
	}
	return fmt.Sprintf("%s, and %s", joinWithCommas(items[:len(items)-1]), items[len(items)-1])
}
