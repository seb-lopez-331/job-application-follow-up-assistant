package followup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/user"
)

type Repository interface {
	Create(ctx context.Context, input CreateFollowUpInput) (*FollowUp, error)
	Read(ctx context.Context, id uuid.UUID) (*FollowUp, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateFollowUpInput) (*FollowUp, error)
	Delete(ctx context.Context, id uuid.UUID) error

	ListDue(ctx context.Context) ([]FollowUp, error)
	ListDueForAllUsers(ctx context.Context) ([]FollowUp, error)
	ListDueForUser(ctx context.Context, userID uuid.UUID) ([]FollowUp, error)
	ListDueBy(ctx context.Context, dueOn time.Time) ([]FollowUp, error)
	ListDueByForAllUsers(ctx context.Context, dueOn time.Time) ([]FollowUp, error)
	ListDueByForUser(ctx context.Context, userID uuid.UUID, dueOn time.Time) ([]FollowUp, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	input CreateFollowUpInput,
) (*FollowUp, error) {
	query := `
		INSERT INTO follow_ups (
			due_on,
			job_application_id,
			message_draft
		) VALUES ($1 $2 $3 $4 $5 $6);`
	var f FollowUp

	err := r.pool.QueryRow(ctx, query,
		input.DueOn,
		input.JobApplicationID,
		input.MessageDraft,
	).Scan(&f)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // 23505 is unique_violation
				return nil, ErrFollowUpDuplicate
			}
			return nil, fmt.Errorf("%w: %v", ErrFollowUpViolation, err)
		}

		return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToInsert, err)
	}

	return &f, nil
}

func (r *PostgresRepository) Read(ctx context.Context, id uuid.UUID) (*FollowUp, error) {
	query := `
		SELECT
			id,
			due_on,
			job_application_id,
			follow_up_status,
			message_draft,
			sent_at,
			completed_at
		FROM follow_ups WHERE id = $1;
	`

	var f FollowUp
	err := r.pool.QueryRow(ctx, query, id).Scan(&f)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFollowUpNotFound // Record not found
		}
		return nil, fmt.Errorf("%w: %v", ErrFollowUpNotRead, err)
	}

	return &f, nil
}

func (r *PostgresRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateFollowUpInput,
) (*FollowUp, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Update("follow_ups").Where(squirrel.Eq{"id": id})

	if input.DueOn != nil {
		builder = builder.Set("due_on", input.DueOn)
	}

	if input.MessageDraft != nil {
		builder = builder.Set("message_draft", input.MessageDraft)
	}

	if input.FollowUpStatus != nil {
		builder = builder.Set("follow_up_status", input.FollowUpStatus)
	}

	if input.SentAt != nil {
		builder = builder.Set("sent_at", input.SentAt)
	}

	if input.CompletedAt != nil {
		builder = builder.Set("completed_at", input.CompletedAt)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrFollowUpNotFound // Record not found
			}
			return nil, fmt.Errorf("%w: %v", ErrFollowUpViolation, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToUpdate, err)
	}

	var f FollowUp
	err = r.pool.QueryRow(ctx, query, args...).Scan(&f)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToUpdate, err)
	}
	return &f, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM follow_ups WHERE id = $1;`

	if _, err := r.pool.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("%w: %v", ErrFollowUpFailedToDelete)
	}

	return nil
}

func (r *PostgresRepository) ListDue(ctx context.Context) ([]FollowUp, error) {
	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	return r.ListDueForUser(ctx, u.ID)
}

func (r *PostgresRepository) ListDueForUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]FollowUp, error) {
	query := `
		SELECT f.id, f.due_on, f.job_application_id, f.follow_up_status,
			f.message_draft, f.sent_at, f.completed_at
		FROM follow_ups AS f
		JOIN job_applications AS ja ON ja.id = f.job_application_id
		WHERE ja.user_id = $1
		AND f.follow_up_status = 'pending'
		ORDER BY f.due_on;
	`
	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	if !u.IsAdmin && u.ID != userID {
		return nil, user.ErrForbidden
	}
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToList, err)
	}
	defer rows.Close()
	return r.listRows(rows)
}

func (r *PostgresRepository) ListDueForAllUsers(ctx context.Context) ([]FollowUp, error) {
	query := `
		SELECT
			id,
			due_on,
			job_application_id,
			follow_up_status,
			message_draft,
			sent_at,
			completed_at
		FROM follow_ups WHERE follow_up_status = 'pending;
	`
	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	if !u.IsAdmin {
		return nil, user.ErrForbidden
	}
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToList, err)
	}
	defer rows.Close()
	return r.listRows(rows)
}

func (r *PostgresRepository) ListDueBy(ctx context.Context, dueOn time.Time) ([]FollowUp, error) {
	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	return r.ListDueByForUser(ctx, u.ID, dueOn)
}

func (r *PostgresRepository) ListDueByForUser(
	ctx context.Context,
	userID uuid.UUID,
	dueOn time.Time,
) ([]FollowUp, error) {
	query := `
		SELECT f.id, f.due_on, f.job_application_id, f.follow_up_status,
			f.message_draft, f.sent_at, f.completed_at
		FROM follow_ups AS f
		JOIN job_applications AS ja ON ja.id = f.job_application_id
		WHERE ja.user_id = $1
		AND f.follow_up_status = 'pending'
		AND f.due_on <= $2
		ORDER BY f.due_on;
	`
	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	if !u.IsAdmin && u.ID != userID {
		return nil, user.ErrForbidden
	}
	rows, err := r.pool.Query(ctx, query, userID, dueOn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToList, err)
	}
	defer rows.Close()
	return r.listRows(rows)
}

func (r *PostgresRepository) ListDueByForAllUsers(
	ctx context.Context,
	dueOn time.Time,
) ([]FollowUp, error) {
	query := `
		SELECT
			id,
			due_on,
			job_application_id,
			follow_up_status,
			message_draft,
			sent_at,
			completed_at
		FROM follow_ups WHERE NOW() > $1;
	`
	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	if !u.IsAdmin {
		return nil, user.ErrForbidden
	}
	rows, err := r.pool.Query(ctx, query, dueOn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToList, err)
	}
	defer rows.Close()
	return r.listRows(rows)
}

func (r *PostgresRepository) listRows(rows pgx.Rows) ([]FollowUp, error) {
	var followUps []FollowUp

	for rows.Next() {
		f := &FollowUp{}
		err := rows.Scan(&f)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToList, err)
		}
		followUps = append(followUps, *f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFollowUpFailedToList, err)
	}

	return followUps, nil
}
