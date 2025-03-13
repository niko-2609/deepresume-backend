package models

type JobPostingRequest struct {
	JobDescription string `json:"jobDescription"`
	JobTitle       string `json:"jobTitle"`
	Company        string `json:"company"`
}

type JobAnalysis struct {
	Keywords       []string `json:"keywords"`
	RequiredSkills []string `json:"requiredSkills"`
	OptionalSkills []string `json:"optionalSkills"`
	Experience     string   `json:"experience"`
	Education      string   `json:"education"`
}
