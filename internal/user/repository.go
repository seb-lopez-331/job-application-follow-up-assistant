package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	Read(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateUserInput struct {
	Username    string
	Email       string
	Password    string
	Preferences *PreferencesInput
	IsAdmin     bool
}

type UpdateUserInput struct {
	Email       *string
	Password    *string
	Preferences *PreferencesInput
	IsAdmin     *bool
}

type PreferencesInput struct {
	FollowUpSchedule *FollowUpSchedule `json:"follow_up_schedule"`
	Mode             *string           `json:"mode"` // light, dark, or system
	ItemsPerPage     *int              `json:"items_per_page"`
	Notifications    *bool             `json:"notifications"`
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (
			username,
			email,
			password_hash,
			preferences,
			is_admin
		)
		VALUES ($1 $2 $3 $4 $5)
	`
	var user User

	err = r.pool.QueryRow(ctx, query,
		input.Username,
		input.Email,
		string(hashedBytes),
		input.Preferences,
		input.IsAdmin,
	).Scan(&user)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("duplicate user entry: %w", err)
			}
			return nil, fmt.Errorf("postgres side error: %w", err)
		}
		return nil, fmt.Errorf("db unable to communicate with client: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepository) Read(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, preferences
		FROM users WHERE id = $1
	`
	var user User
	err := r.pool.QueryRow(ctx, query, id).Scan(&user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("user not found: %w", err)
			}
			return nil, fmt.Errorf("postgres side error: %w", err)
		}
		return nil, fmt.Errorf("db unable to communicate with client: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateUserInput,
) (*User, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Update("users").Where(squirrel.Eq{"id": id})

	if input.Email != nil {
		builder = builder.Set("email", input.Email)
	}
	if input.Password != nil {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		builder = builder.Set("hashed_bytes", string(hashedBytes))
	}
	if input.Preferences != nil {
		builder = builder.Set("email", input.Preferences)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("user not found: %w", err)
			}
			return nil, fmt.Errorf("postgres side error: %w", err)
		}
		return nil, fmt.Errorf("db unable to communicate with client: %w", err)
	}

	var user User
	err = r.pool.QueryRow(ctx, query, args...).Scan(&user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("user not found: %w", err)
			}
			return nil, fmt.Errorf("postgres side error: %w", err)
		}
		return nil, fmt.Errorf("db unable to communicate with client: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("user not found: %w", err)
			}
			return fmt.Errorf("postgres side error: %w", err)
		}
		return fmt.Errorf("db unable to communicate with client: %w", err)
	}

	return nil
}
