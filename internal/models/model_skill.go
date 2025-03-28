package models

import (
	"time"
)

// Skill represents a general skill in the system
type Skill struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserSkill represents the relationship between users and skills
type UserSkill struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"userId" gorm:"not null"`
	SkillID     uint      `json:"skillId" gorm:"not null"`
	Proficiency string    `json:"proficiency"` // e.g., "Beginner", "Intermediate", "Expert"
	YearsOfExp  float32   `json:"yearsOfExp"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	User  User  `json:"-" gorm:"foreignKey:UserID"`
	Skill Skill `json:"skill" gorm:"foreignKey:SkillID"`
}
