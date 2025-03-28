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

// type Education struct {
// 	Institution string     `json:"institution"`
// 	Degree      string     `json:"degree"`
// 	Field       string     `json:"field"`
// 	StartDate   time.Time  `json:"startDate"`
// 	EndDate     *time.Time `json:"endDate,omitempty"`
// 	GPA         float64    `json:"gpa,omitempty"`
// }

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

// Resume represents a specific version of a user's resume
type Resume struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"userId" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`     // e.g., "Software Engineer - Google"
	Description string    `json:"description"`              // Optional description
	JobTitle    string    `json:"jobTitle"`                 // Target job title
	Company     string    `json:"company"`                  // Target company
	Content     string    `json:"content" gorm:"type:text"` // JSON content of the resume
	IsDefault   bool      `json:"isDefault" gorm:"default:false"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}

// ResumeSkill represents skills included in a specific resume version
type ResumeSkill struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ResumeID  uint      `json:"resumeId" gorm:"not null"`
	SkillID   uint      `json:"skillId" gorm:"not null"`
	Order     int       `json:"order"` // For custom ordering of skills in the resume
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Resume Resume `json:"-" gorm:"foreignKey:ResumeID"`
	Skill  Skill  `json:"skill" gorm:"foreignKey:SkillID"`
}
