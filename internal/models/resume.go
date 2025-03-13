package models

import "time"

type PersonalInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
	GitHub   string `json:"github,omitempty"`
}

type Experience struct {
	Company     string     `json:"company"`
	Title       string     `json:"title"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Current     bool       `json:"current"`
	Location    string     `json:"location"`
	Description []string   `json:"description"`
}

type Education struct {
	Institution string     `json:"institution"`
	Degree      string     `json:"degree"`
	Field       string     `json:"field"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	GPA         float64    `json:"gpa,omitempty"`
}

type Project struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	URL          string   `json:"url,omitempty"`
}

type ResumeContent struct {
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Summary      string       `json:"summary"`
	Experience   []Experience `json:"experience"`
	Education    []Education  `json:"education"`
	Skills       []string     `json:"skills"`
	Projects     []Project    `json:"projects"`
}

type GenerateResumeRequest struct {
	PersonalInfo   PersonalInfo   `json:"personalInfo"`
	JobDescription string         `json:"jobDescription"`
	ExistingResume *ResumeContent `json:"existingResume,omitempty"`
}

type GenerateResumeResponse struct {
	ResumeContent ResumeContent `json:"resumeContent"`
	Suggestions   []string      `json:"suggestions"`
}
