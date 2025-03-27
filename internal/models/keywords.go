package models

// Keyword represents a single keyword and its metadata
type Keyword struct {
	Word  string  `json:"word"`
	Score float64 `json:"score"`
	Count int     `json:"count"`
}

// KeywordVectorResult represents a document with keywords and similarity score
type KeywordVectorResult struct {
	ID         int       `json:"id"`
	SourceType string    `json:"source_type"`
	SourceID   string    `json:"source_id"`
	Keywords   []Keyword `json:"keywords"`
	Similarity float64   `json:"similarity"`
}
