package jobapplication

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Repository interface {
	Create(ctx context.Context, app JobApplication) (JobApplication, error)
}

type PostgresRepository struct {

}

func (r *PostgresRepository) Create(
	ctx context.Context,
	app JobApplication,
) (*JobApplication, error) {
	// Repository passes the context into the database operation
	// The actual SQL an returned columns go in here
	return r.insertJobApplication(ctx, app)
}

func (r *PostgresRepository) insertJobApplication(
	ctx context.Context,
	app JobApplication,
) (*JobApplication, error) {

}