package models

import (
	"time"
)

type User struct {
	ID             uint             `json:"id" gorm:"primaryKey"`
	Email          string           `json:"email" gorm:"unique;not null"`
	FullName       string           `json:"fullName" gorm:"not null"`
	Phone          string           `json:"phone"`
	Location       string           `json:"location"`
	Title          string           `json:"title"`
	Summary        string           `json:"summary"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	WorkExperience []WorkExperience `json:"workExperience" gorm:"foreignKey:UserID"`
	Education      []Education      `json:"education" gorm:"foreignKey:UserID"`
}
