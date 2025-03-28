package repository

import (
	"context"
	"errors"
	"time"

	"github.com/nikolai/ai-resume-builder/backend/internal/interfaces"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
)

type UserRepository struct {
	db interfaces.DB
}

func NewUserRepository(db interfaces.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.CreateUserTx(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

// CreateUserTx creates a new user within a transaction
func (r *UserRepository) CreateUserTx(ctx context.Context, tx interfaces.Tx, user *models.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Use a more efficient insert by specifying the fields
	result := tx.Model(user).Select("email", "full_name", "phone", "location", "title", "summary", "created_at", "updated_at").Create(user)
	return result.Error
}

// CreateWorkExperiencesTx creates multiple work experiences within a transaction
func (r *UserRepository) CreateWorkExperiencesTx(ctx context.Context, tx interfaces.Tx, experiences []models.WorkExperience) error {
	if len(experiences) == 0 {
		return nil
	}
	result := tx.CreateInBatches(experiences, 100)
	return result.Error
}

// CreateEducationsTx creates multiple education records within a transaction
func (r *UserRepository) CreateEducationsTx(ctx context.Context, tx interfaces.Tx, educations []models.Education) error {
	if len(educations) == 0 {
		return nil
	}
	result := tx.CreateInBatches(educations, 100)
	return result.Error
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Select("users.id, users.email, users.full_name, users.phone, users.location, users.title, users.summary, users.created_at, users.updated_at").
		Where("users.id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Select("users.id, users.email, users.full_name, users.phone, users.location, users.title, users.summary, users.created_at, users.updated_at").
		Where("users.email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
