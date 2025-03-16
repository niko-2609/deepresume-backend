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
	query := `
		INSERT INTO users (email, full_name, phone, location, title, summary, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return tx.QueryRowContext(ctx, query,
		user.Email,
		user.FullName,
		user.Phone,
		user.Location,
		user.Title,
		user.Summary,
		now,
	).Scan(&user.ID)
}

// CreateWorkExperiencesTx creates multiple work experiences within a transaction
func (r *UserRepository) CreateWorkExperiencesTx(ctx context.Context, tx interfaces.Tx, experiences []models.WorkExperience) error {
	query := `
		INSERT INTO work_experiences (user_id, company, title, location, start_date, end_date, current, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, exp := range experiences {
		_, err := tx.ExecContext(ctx, query,
			exp.UserID,
			exp.Company,
			exp.Title,
			exp.Location,
			exp.StartDate,
			exp.EndDate,
			exp.IsCurrent,
			exp.Description,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateEducationsTx creates multiple education records within a transaction
func (r *UserRepository) CreateEducationsTx(ctx context.Context, tx interfaces.Tx, educations []models.Education) error {
	query := `
		INSERT INTO educations (user_id, school, degree, field, location, start_date, end_date, current, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	for _, edu := range educations {
		_, err := tx.ExecContext(ctx, query,
			edu.UserID,
			edu.School,
			edu.Degree,
			edu.Field,
			edu.Location,
			edu.StartDate,
			edu.EndDate,
			edu.IsCurrent,
			edu.Description,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	query := `
		SELECT id, email, full_name, phone, location, title, summary, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Phone,
		&user.Location,
		&user.Title,
		&user.Summary,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, full_name, phone, location, title, summary, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Phone,
		&user.Location,
		&user.Title,
		&user.Summary,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE users
		SET full_name = $1, phone = $2, location = $3, title = $4, summary = $5, updated_at = $6
		WHERE id = $7`

	user.UpdatedAt = time.Now()
	result, err := tx.ExecContext(ctx, query,
		user.FullName,
		user.Phone,
		user.Location,
		user.Title,
		user.Summary,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return tx.Commit()
}
