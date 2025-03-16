package service

import (
	"context"
)

// CreateUserOnboarding handles the creation of user profile, work experience, and education in a transaction
func (s *Service) CreateUserOnboarding(ctx context.Context, data OnboardingData) error {
	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	if err := tx.Create(&data.User).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add user ID to work experiences
	for i := range data.WorkExperience {
		data.WorkExperience[i].UserID = data.User.ID
	}

	// Create work experiences
	if len(data.WorkExperience) > 0 {
		if err := tx.Create(&data.WorkExperience).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Add user ID to education records
	for i := range data.Education {
		data.Education[i].UserID = data.User.ID
	}

	// Create education records
	if len(data.Education) > 0 {
		if err := tx.Create(&data.Education).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}
