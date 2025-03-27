package service

import (
	"context"
	"fmt"

	"github.com/nikolai/ai-resume-builder/backend/internal/interfaces"
)

type ResumeService struct {
	db             interfaces.DB
	keywordService *KeywordService
	llmService     *LLMService
}

func NewResumeService(db interfaces.DB, keywordService *KeywordService, llmService *LLMService) *ResumeService {
	return &ResumeService{
		db:             db,
		keywordService: keywordService,
		llmService:     llmService,
	}
}

// ResumeStreamHandler is a function that handles streaming resume chunks
type ResumeStreamHandler func(chunk string, done bool) error

// GenerateResume generates a resume based on the job description
func (s *ResumeService) GenerateResume(ctx context.Context, jobDescription string) (string, error) {
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

	// If LLM service is available, use it to generate the resume
	if s.llmService != nil {
		// Prepare a prompt for the LLM
		prompt := s.buildPrompt(keywordStrings, jobDescription)

		// Call the LLM service to generate the resume
		resume, err := s.llmService.GenerateContent(ctx, "deepseek-coder", prompt)
		if err != nil {
			return "", fmt.Errorf("failed to generate resume with LLM: %v", err)
		}
		return resume, nil
	}

	return "", fmt.Errorf("LLM service is not available")
}

// StreamGenerateResume generates a resume and streams the results to the handler
func (s *ResumeService) StreamGenerateResume(ctx context.Context, jobDescription string, handler ResumeStreamHandler) error {
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
		prompt := s.buildPrompt(keywordStrings, jobDescription)

		// Stream the LLM responses
		err := s.llmService.StreamGenerateContent(ctx, "deepseek-r1", prompt, func(chunk string, done bool) error {
			return handler(chunk, done)
		})

		if err != nil {
			return fmt.Errorf("failed to stream resume generation: %v", err)
		}

		return nil
	}

	return fmt.Errorf("LLM service is not available")
}

// buildPrompt creates a prompt for the LLM based on the extracted keywords
func (s *ResumeService) buildPrompt(keywordStrings []string, jobDescription string) string {
	return fmt.Sprintf(`
 Generate an ATS-optimized resume tailored to the following job description. 
 Ensure the resume includes relevant keywords, a professional format, and highlights key skills, experience, and achievements
Job Description: 
%s

The resume should include:
1. A professional summary highlighting expertise in skills mentioned in the job description
2. A skills section that emphasizes the keywords
3. Experience section with relevant bullet points
4. Education section
5. Any other relevant sections

Format the resume in Markdown.
`, jobDescription)
}
