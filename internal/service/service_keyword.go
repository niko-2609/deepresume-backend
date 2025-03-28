package service

import (
	"sort"
	"strings"

	"github.com/nikolai/ai-resume-builder/backend/internal/interfaces"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
)

type KeywordService struct {
	db interfaces.DB
}

func NewKeywordService(db interfaces.DB) *KeywordService {
	return &KeywordService{
		db: db,
	}
}

// Common English stopwords
var stopwords = map[string]bool{
	"a": true, "about": true, "above": true, "after": true, "again": true, "against": true,
	"all": true, "am": true, "an": true, "and": true, "any": true, "are": true, "as": true,
	"at": true, "be": true, "because": true, "been": true, "before": true, "being": true,
	"below": true, "between": true, "both": true, "but": true, "by": true, "could": true,
	"did": true, "do": true, "does": true, "doing": true, "down": true, "during": true,
	// ... other stopwords (keep the ones you had before)
}

// Extract and Rank Keywords from text
func (s *KeywordService) ExtractAndRankKeywords(text string) []models.Keyword {
	// Preprocessing: lowercase the text and replace punctuation with spaces
	text = strings.ToLower(text)
	text = cleanText(text)

	// Define common multi-word technical phrases to look for
	commonPhrases := []string{
		// Job titles
		"software engineer", "backend engineer", "frontend engineer", "full stack engineer",
		"software developer", "web developer", "devops engineer", "data engineer",
		"system administrator", "database administrator", "cloud architect", "solutions architect",
		"product manager", "scrum master", "technical lead", "engineering manager",

		// Technologies and frameworks
		"restful api", "restful apis", "microservices architecture", "service oriented architecture",
		"continuous integration", "continuous deployment", "ci/cd pipeline", "git workflow",
		"test driven development", "agile methodology", "scrum methodology", "kanban methodology",

		// Skills and concepts
		"cloud infrastructure", "distributed systems", "system design", "database design",
		"api design", "object oriented programming", "functional programming", "version control",
		"data structures", "design patterns", "unit testing", "integration testing",

		// Cloud platforms and tools
		"aws cloud", "microsoft azure", "google cloud", "cloud computing",
		"amazon web services", "aws lambda", "aws ec2", "aws s3",
		"docker containers", "kubernetes orchestration", "terraform", "infrastructure as code",

		// Programming languages with context
		"golang development", "python programming", "javascript framework", "typescript development",
		"java enterprise", "c++ programming", "react development", "node.js development",

		// Databases
		"postgresql database", "mysql database", "mongodb database", "redis cache",
		"elasticsearch", "sqlite database", "dynamodb", "database optimization",

		// Machine learning and data
		"machine learning", "data analysis", "data visualization", "big data",
		"natural language processing", "computer vision", "predictive modeling", "neural networks",
	}

	// Count phrases first
	phraseCount := make(map[string]int)
	for _, phrase := range commonPhrases {
		paddedPhrase := " " + phrase + " "
		paddedText := " " + text + " "
		count := strings.Count(paddedText, paddedPhrase)
		if count > 0 {
			phraseCount[phrase] = count
			text = strings.ReplaceAll(paddedText, paddedPhrase, " "+strings.Repeat("X", len(phrase))+" ")
			text = strings.TrimSpace(text)
		}
	}

	// Split text into tokens for further processing
	text = cleanText(text)
	tokens := strings.Fields(text)

	// Extract potential n-gram phrases
	if len(tokens) >= 2 {
		for i := 0; i < len(tokens)-1; i++ {
			word1 := tokens[i]
			word2 := tokens[i+1]

			if stopwords[word1] || stopwords[word2] ||
				len(word1) <= 2 || len(word2) <= 2 ||
				strings.Contains(word1, "X") || strings.Contains(word2, "X") {
				continue
			}

			bigram := word1 + " " + word2

			if _, exists := phraseCount[bigram]; !exists {
				paddedBigram := " " + bigram + " "
				paddedText := " " + text + " "
				count := strings.Count(paddedText, paddedBigram)
				if count > 0 {
					phraseCount[bigram] = count
					text = strings.ReplaceAll(paddedText, paddedBigram, " "+strings.Repeat("X", len(bigram))+" ")
					text = strings.TrimSpace(text)
				}
			}
		}
	}

	// Process individual words
	words := strings.Fields(text)
	wordCount := make(map[string]int)
	totalTerms := 0

	// Add phrase counts
	for phrase, count := range phraseCount {
		wordCount[phrase] = count
		totalTerms += count
	}

	// Add individual word counts
	for _, word := range words {
		if strings.Contains(word, "X") {
			continue
		}
		if len(word) > 2 && !stopwords[word] {
			wordCount[word]++
			totalTerms++
		}
	}

	// Calculate term frequency and create keywords
	var keywords []models.Keyword
	for term, count := range wordCount {
		tf := float64(count) / float64(totalTerms)
		wordWeight := 1.0
		if strings.Contains(term, " ") {
			wordWeight = 1.5
		}
		score := tf * wordWeight

		keywords = append(keywords, models.Keyword{
			Word:  term,
			Score: score,
			Count: count,
		})
	}

	// Sort by score
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Score > keywords[j].Score
	})

	// Return top keywords
	maxKeywords := 50
	if len(keywords) < maxKeywords {
		maxKeywords = len(keywords)
	}
	return keywords[:maxKeywords]
}

// cleanText replaces punctuation with spaces and cleans up the text
func cleanText(text string) string {
	replacer := strings.NewReplacer(
		".", " ", ",", " ", ";", " ", ":", " ",
		"!", " ", "?", " ", "(", " ", ")", " ",
		"[", " ", "]", " ", "{", " ", "}", " ",
		"/", " ", "\\", " ", "-", " ", "_", " ",
		"+", " ", "=", " ", "*", " ", "&", " ",
		"%", " ", "$", " ", "#", " ", "@", " ",
	)
	text = replacer.Replace(text)

	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return strings.TrimSpace(text)
}
