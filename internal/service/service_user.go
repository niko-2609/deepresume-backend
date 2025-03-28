package service

import (
	"context"
	"errors"

	"github.com/nikolai/ai-resume-builder/backend/internal/interfaces"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
	"github.com/nikolai/ai-resume-builder/backend/internal/repository"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepository
	db       interfaces.DB
}

type OnboardingData struct {
	User           models.User             `json:"user"`
	WorkExperience []models.WorkExperience `json:"workExperience"`
	Education      []models.Education      `json:"education"`
}

func NewUserService(userRepo *repository.UserRepository, db interfaces.DB) *UserService {
	return &UserService{
		userRepo: userRepo,
		db:       db,
	}
}

// CreateUser creates a new user with onboarding information
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	// Validate required fields
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.FullName == "" {
		return errors.New("full name is required")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Create user
	return s.userRepo.CreateUser(ctx, user)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	// Validate required fields
	if user.ID == 0 {
		return errors.New("user ID is required")
	}
	if user.FullName == "" {
		return errors.New("full name is required")
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return s.userRepo.UpdateUser(ctx, user)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.userRepo.GetUserByEmail(ctx, email)
}

// CreateUserOnboarding handles the creation of user profile, work experience, and education in a transaction
func (s *UserService) CreateUserOnboarding(ctx context.Context, data OnboardingData) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create user
	if err := s.userRepo.CreateUserTx(ctx, tx, &data.User); err != nil {
		return err
	}

	// Add user ID to work experiences
	for i := range data.WorkExperience {
		data.WorkExperience[i].UserID = data.User.ID
	}

	// Create work experiences
	if len(data.WorkExperience) > 0 {
		if err := s.userRepo.CreateWorkExperiencesTx(ctx, tx, data.WorkExperience); err != nil {
			return err
		}
	}

	// Add user ID to education records
	for i := range data.Education {
		data.Education[i].UserID = data.User.ID
	}

	// Create education records
	if len(data.Education) > 0 {
		if err := s.userRepo.CreateEducationsTx(ctx, tx, data.Education); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetUserWithDetails retrieves a user by ID with their work experience and education
func (s *UserService) GetUserWithDetails(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	// Use a single query with proper indexing
	err := s.db.WithContext(ctx).
		Select("users.id, users.full_name, users.email, users.phone, users.location, users.title, users.summary").
		Joins("LEFT JOIN work_experiences ON users.id = work_experiences.user_id").
		Joins("LEFT JOIN educations ON users.id = educations.user_id").
		Preload("WorkExperience", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, user_id, company, title, location, start_date, end_date, is_current, description").
				Order("start_date DESC")
		}).
		Preload("Education", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, user_id, school, degree, field, location, start_date, end_date, is_current, description").
				Order("start_date DESC")
		}).
		Where("users.id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
