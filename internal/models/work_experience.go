package models

import (
	"time"
)

type WorkExperience struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"userId" gorm:"not null"`
	Company     string     `json:"company" gorm:"not null"`
	Title       string     `json:"title" gorm:"not null"`
	Location    string     `json:"location"`
	StartDate   time.Time  `json:"startDate" gorm:"not null"`
	EndDate     *time.Time `json:"endDate"`
	IsCurrent   bool       `json:"isCurrent" gorm:"default:false"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}
