package jobapplication

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ColID          = "id"
	ColUserID      = "user_id"
	ColAppliedOn   = "applied_on"
	ColCompanyName = "company_name"
	ColRoleTitle   = "role_title"

	ColDatePosted      = "date_posted"
	ColDescription     = "desc"
	ColJobURL          = "job_url"
	ColNotes           = "notes"
	ColLocation        = "loc"
	ColPostedAt        = "posted_at"
	ColWorkArrangement = "work_arrangement"
	ColStage           = "stage"
	ColNextFollowUpOn  = "next_follow_up_on"
)

type Repository interface {
	Create(ctx context.Context, input CreateJobApplicationInput) (*JobApplication, error)
	Read(ctx context.Context, id uuid.UUID) (*JobApplication, error)
	Update(
		ctx context.Context,
		id uuid.UUID,
		input UpdateJobApplicationInput,
	) (*JobApplication, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter JobApplicationFilter) ([]JobApplication, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	input CreateJobApplicationInput,
) (*JobApplication, error) {
	query := `
		INSERT INTO job_applications (
			user_id,
			applied_on
			company_name,
			role_title,

			date_posted,
			descr,
			job_url,
			notes,
			next_follow_up_on,
			loc,
			posted_at,
			stage,
			work_arrangement
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	var app JobApplication

	err := r.pool.QueryRow(ctx, query,
		input.UserID,
		input.AppliedOn,
		input.CompanyName,
		input.RoleTitle,

		input.DatePosted,
		input.Description,
		input.JobURL,
		input.Notes,
		input.NextFollowUpOn,
		input.Location,
		input.PostedAt,
		input.Stage,
		input.WorkArrangement,
	).Scan(&app)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // 23505 is unique_violation
				return nil, ErrJobApplicationDuplicate
			}
			return nil, fmt.Errorf("%w: %v", ErrJobApplicationInvalidInput, err)
		}

		return nil, fmt.Errorf("%w: %v", ErrJobApplicationFailedToInsert, err)
	}

	return &app, nil
}

func (r *PostgresRepository) Read(ctx context.Context, id uuid.UUID) (*JobApplication, error) {
	query := `
		SELECT
			id,
			user_id,
			applied_on,
			company_name,
			role_title,
			date_posted,
			descr,
			job_url,
			notes,
			next_follow_up_on,
			loc,
			posted_at,
			work_arrangement,
			stage,
			next_follow_up_on
		FROM job_applications WHERE id = $1
	`

	var app JobApplication
	err := r.pool.QueryRow(ctx, query, id).Scan(&app)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobApplicationNotFound // Record not found
		}
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationNotRead, err)
	}

	return &app, nil
}

func (r *PostgresRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateJobApplicationInput,
) (*JobApplication, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Update("job_applications").Where(squirrel.Eq{"id": id})

	// 1. Handle standard Pointers (skip if nil)
	if input.AppliedOn != nil {
		builder = builder.Set("applied_on", input.AppliedOn)
	}
	if input.CompanyName != nil {
		builder = builder.Set("company_name", input.CompanyName)
	}
	if input.RoleTitle != nil {
		builder = builder.Set("role_title", input.RoleTitle)
	}
	if input.DatePosted != nil {
		builder = builder.Set("date_posted", input.DatePosted)
	}
	if input.Description != nil {
		builder = builder.Set("description", input.Description)
	}
	if input.Notes != nil {
		builder = builder.Set("notes", input.Notes)
	}
	if input.NextFollowUpOn != nil {
		builder = builder.Set("next_follow_up_on", input.NextFollowUpOn)
	}
	if input.Location != nil {
		builder = builder.Set("location", input.Location)
	}
	if input.PostedAt != nil {
		builder = builder.Set("posted_at", input.PostedAt)
	}
	if input.Stage != nil {
		builder = builder.Set("stage", input.Stage)
	}
	if input.WorkArrangement != nil {
		builder = builder.Set("work_arrangement", input.WorkArrangement)
	}

	// 2. Handle Custom Pointer Types (e.g., net/url or sub-structs)
	if input.JobURL != nil {
		// Convert *url.URL to string for the DB storage
		builder = builder.Set("job_url", input.JobURL.String())
	}
	if input.UserID != nil {
		builder = builder.Set("user_id", input.UserID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobApplicationNotFound // Record not found
		}
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationFailedToInsert, err)
	}

	var app JobApplication
	err = r.pool.QueryRow(ctx, query, args...).Scan(&app)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationFailedToInsert, err)
	}
	return &app, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM job_applications WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrJobApplicationFailedToDelete, err)
	}

	return nil
}

func (r *PostgresRepository) List(
	ctx context.Context,
	filter JobApplicationFilter,
) ([]JobApplication, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Select(
		ColUserID,
		ColAppliedOn,
		ColCompanyName,
		ColRoleTitle,
		ColDatePosted,
		ColDescription,
		ColJobURL,
		ColNotes,
		ColLocation,
		ColPostedAt,
		ColWorkArrangement,
		ColStage,
		ColNextFollowUpOn,
	).From("job_applications")

	// 2. Dynamically add WHERE clauses based on what the user provided
	// Squirrel will automatically combine these with "AND"
	if filter.ID != nil {
		builder = builder.Where(squirrel.Eq{ColID: filter.ID})
	}

	if filter.UserID != nil {
		builder = builder.Where(squirrel.Eq{ColUserID: filter.UserID})
	}

	if filter.AppliedOn != nil {
		builder = builder.Where(squirrel.Eq{ColAppliedOn: filter.AppliedOn})
	} else {
		if filter.AppliedBefore != nil {
			builder = builder.Where(squirrel.Lt{ColAppliedOn: filter.AppliedBefore})
		}
		if filter.AppliedAfter != nil {
			builder = builder.Where(squirrel.Gt{ColAppliedOn: filter.AppliedAfter})
		}
	}

	if filter.CompanyName != nil {
		builder = builder.Where(squirrel.Eq{ColCompanyName: filter.CompanyName})
	}

	if filter.RoleTitle != nil {
		builder = builder.Where(squirrel.Eq{ColRoleTitle: filter.RoleTitle})
	}

	if filter.DatePosted != nil {
		builder = builder.Where(squirrel.Eq{ColDatePosted: filter.DatePosted})
	} else {
		if filter.PostedBefore != nil {
			builder = builder.Where(squirrel.Lt{ColDatePosted: filter.PostedBefore})
		}
		if filter.PostedAfter != nil {
			builder = builder.Where(squirrel.Gt{ColDatePosted: filter.PostedAfter})
		}
	}

	if filter.Description != nil {
		builder = builder.Where(squirrel.Eq{ColDescription: filter.Description})
	}

	if filter.JobURL != nil {
		builder = builder.Where(squirrel.Eq{ColJobURL: filter.JobURL.String()})
	}

	if filter.Notes != nil {
		builder = builder.Where(squirrel.Eq{ColNotes: filter.Notes})
	}

	if filter.Location != nil {
		builder = builder.Where(squirrel.Eq{ColLocation: filter.Location})
	}

	if filter.PostedAt != nil {
		builder = builder.Where(squirrel.Eq{ColPostedAt: filter.PostedAt})
	}

	if filter.WorkArrangement != nil {
		builder = builder.Where(squirrel.Eq{ColWorkArrangement: filter.WorkArrangement})
	}

	if filter.Stage != nil {
		builder = builder.Where(squirrel.Eq{ColStage: filter.Stage})
	}

	if filter.NextFollowUpOn != nil {
		builder = builder.Where(squirrel.Eq{ColNextFollowUpOn: filter.NextFollowUpOn})
	}

	// 3. Compile the final SQL string and its arguments
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationFailedToList, err)
	}

	// 4. Execute it safely (Squirrel generated the correct $1, $2, etc.)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationFailedToList, err)
	}
	defer rows.Close()

	var apps []JobApplication

	for rows.Next() {
		app := &JobApplication{}
		var rawURL *string
		err := rows.Scan(
			&app.ID,
			&app.UserID,
			&app.AppliedOn,
			&app.CompanyName,
			&app.RoleTitle,
			&app.DatePosted,
			&app.Description,
			&rawURL, // Scan the text column into our string pointer
			&app.Notes,
			&app.Location,
			&app.PostedAt,
			&app.WorkArrangement,
			&app.Stage,
			&app.NextFollowUpOn,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrJobApplicationFailedToList, err)
		}

		if rawURL != nil {
			parsedURL, err := url.Parse(*rawURL)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrJobApplicationFailedToList, err)
			}
			app.JobURL = parsedURL
		}

		apps = append(apps, *app)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationFailedToList, err)
	}

	return apps, nil
}
