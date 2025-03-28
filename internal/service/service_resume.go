package service

import (
	"context"
	"fmt"

	"github.com/nikolai/ai-resume-builder/backend/internal/interfaces"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
)

type ResumeService struct {
	db             interfaces.DB
	keywordService *KeywordService
	llmService     *LLMService
	userService    *UserService
}

func NewResumeService(db interfaces.DB, keywordService *KeywordService, llmService *LLMService, userService *UserService) *ResumeService {
	return &ResumeService{
		db:             db,
		keywordService: keywordService,
		llmService:     llmService,
		userService:    userService,
	}
}

// ResumeStreamHandler is a function that handles streaming resume chunks
type ResumeStreamHandler func(chunk string, done bool) error

// GenerateResume generates a resume based on the job description and streams the results
func (s *ResumeService) GenerateResume(ctx context.Context, userID uint, jobDescription string, handler ResumeStreamHandler) error {
	// Fetch user with related data
	user, err := s.userService.GetUserWithDetails(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user data: %v", err)
	}

	// Extract keywords from job description
	keywords := s.keywordService.ExtractAndRankKeywords(jobDescription)

	// Prepare keywords for the LLM
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

	// If LLM service is available, use it to generate the resume with streaming
	if s.llmService != nil {
		// Prepare a prompt for the LLM
		prompt := s.buildPrompt(keywordStrings, jobDescription, user)

		// Stream the LLM responses
		err := s.llmService.StreamGenerateContent(ctx, "cusmodel1.2", prompt, func(chunk string, done bool) error {
			return handler(chunk, done)
		})

		if err != nil {
			return fmt.Errorf("failed to stream resume generation: %v", err)
		}

		return nil
	}

	return fmt.Errorf("LLM service is not available")
}

// buildPrompt creates a prompt for the LLM based on the extracted keywords and user data
func (s *ResumeService) buildPrompt(keywordStrings []string, jobDescription string, user *models.User) string {
	// Format personal information
	personalInfo := fmt.Sprintf("Name: %s\nEmail: %s\nPhone: %s\nLocation: %s\nTitle: %s\nSummary: %s\n\n",
		user.FullName,
		user.Email,
		user.Phone,
		user.Location,
		user.Title,
		user.Summary,
	)

	// Format work experience
	experience := ""
	for _, exp := range user.WorkExperience {
		experience += fmt.Sprintf("- %s at %s (%s - %s)\n", exp.Title, exp.Company, exp.StartDate.Format("Jan 2006"), getEndDate(exp))
		experience += fmt.Sprintf("  Location: %s\n", exp.Location)
		experience += fmt.Sprintf("  Description: %s\n\n", exp.Description)
	}

	// Format education
	education := ""
	for _, edu := range user.Education {
		education += fmt.Sprintf("- %s in %s from %s (%s - %s)\n", edu.Degree, edu.Field, edu.School, edu.StartDate.Format("Jan 2006"), getEducationEndDate(edu))
		education += fmt.Sprintf("  Location: %s\n", edu.Location)
		education += fmt.Sprintf("  Description: %s\n\n", edu.Description)
	}

	// Format skills
	skills := ""
	for _, skill := range keywordStrings {
		skills += fmt.Sprintf("- %s\n", skill)
	}

	return fmt.Sprintf(
		`You are a professional resume writer. Your task is to create an ATS-optimized resume in markdown format using ONLY the information provided below. Do not make up or add any information that is not explicitly provided.

IMPORTANT: Use the exact name, contact details, and information provided in the Personal Information section. Do not modify or change any of these details.

Use markdown syntax for formatting (e.g., # for headings, * for emphasis, etc.).

Personal Information:
%s

Job Description:
%s

Experience:
%s

Skills:
%s

Education:
%s

Generate a professional resume that highlights the candidate's experience and skills in relation to the job description. Use only the information provided above.`,
		personalInfo, jobDescription, experience, skills, education,
	)
}

func getEndDate(exp models.WorkExperience) string {
	if exp.IsCurrent {
		return "Present"
	}
	if exp.EndDate != nil {
		return exp.EndDate.Format("Jan 2006")
	}
	return "Present"
}

func getEducationEndDate(edu models.Education) string {
	if edu.IsCurrent {
		return "Present"
	}
	if edu.EndDate != nil {
		return edu.EndDate.Format("Jan 2006")
	}
	return "Present"
}
