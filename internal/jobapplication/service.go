package jobapplication

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/user"
)

// The service depends on the Repository interface, not directly on PostgreSQL.
// That makes it possible to unit-test scheduling and validation rules with
// a fake repository
type JobApplicationService interface {
	Create(
		ctx context.Context,
		input CreateJobApplicationInput,
	) (*JobApplication, error)

	Read(ctx context.Context, id uuid.UUID) (*JobApplication, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		input UpdateJobApplicationInput,
	) (*JobApplication, error)

	List(
		ctx context.Context,
		filter JobApplicationFilter,
	) ([]JobApplication, error)

	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateJobApplicationInput struct {
	AppliedOn   time.Time
	CompanyName string
	RoleTitle   string

	DatePosted      *time.Time
	Description     *string
	JobURL          *url.URL
	Notes           *string
	NextFollowUpOn  *time.Time
	Location        *string
	PostedAt        *time.Time
	Stage           *string
	User            *user.User
	WorkArrangement *string
}

type UpdateJobApplicationInput struct {
	AppliedOn       *time.Time
	CompanyName     *string
	RoleTitle       *string
	DatePosted      *time.Time
	Description     *string
	JobURL          *url.URL
	Notes           *string
	NextFollowUpOn  *time.Time
	Location        *string
	PostedAt        *time.Time
	Stage           *string
	User            *user.User
	WorkArrangement *string
}

type JobApplicationFilter struct {
	ID              *uuid.UUID
	UserID          *uuid.UUID
	AppliedOn       *time.Time
	CompanyName     *string
	RoleTitle       *string
	DatePosted      *time.Time
	Description     *string
	JobURL          *url.URL
	Notes           *string
	Location        *string
	PostedAt        *time.Time
	WorkArrangement *string
	Stage           *string
	NextFollowUpOn  *time.Time
}

type jobApplicationService struct {
	repo Repository
}

func NewJobApplicationService(repo Repository) *jobApplicationService {
	return &jobApplicationService{repo: repo}
}

func (s *jobApplicationService) Create(
	ctx context.Context,
	input CreateJobApplicationInput,
) (*JobApplication, error) {
	// Validate job application input
	app, err := ToJobApplication(input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationInvalidInput, err)
	}
	if !app.IsValid() {
		return nil, ErrJobApplicationInvalidInput
	}

	if *input.Stage == "applied" {
		due := addBusinessDays(input.AppliedOn, 7)
		input.NextFollowUpOn = &due
	}
	return s.repo.Create(ctx, input)
}

func (s *jobApplicationService) Read(
	ctx context.Context,
	id uuid.UUID,
) (*JobApplication, error) {
	app, err := s.repo.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}

	if !u.IsAdmin && app.UserID != u.ID {
		return nil, user.ErrForbidden
	}

	return app, nil
}

func (s *jobApplicationService) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateJobApplicationInput,
) (*JobApplication, error) {
	// Ensure that the UpdateJobApplicationInput has at least one not null field
	if input.IsEmpty() {
		return nil, ErrJobApplicationInvalidInput
	}
	app, err := ToJobApplication(input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationInvalidInput, err)
	}
	if !app.IsValid() {
		return nil, ErrJobApplicationInvalidInput
	}
	return s.repo.Update(ctx, id, input)
}

func (s *jobApplicationService) Delete(ctx context.Context, id uuid.UUID) error {
	u, ok := user.FromContext(ctx)
	if !ok {
		return user.ErrUserNotFoundInCtx
	}

	if !u.IsAdmin {
		app, err := s.Read(ctx, id)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrJobApplicationFailedToDelete, err)
		}
		if app.UserID != u.ID {
			return user.ErrForbidden
		}
	}

	return s.repo.Delete(ctx, id)
}

func (s *jobApplicationService) List(
	ctx context.Context,
	filter JobApplicationFilter,
) ([]JobApplication, error) {
	// Validate the filter
	app, err := ToJobApplication(filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJobApplicationInvalidInput, err)
	}
	if !app.IsValid() {
		return nil, ErrJobApplicationInvalidInput
	}

	return s.repo.List(ctx, filter)
}

func addBusinessDays(start time.Time, days int) time.Time {
	current := start
	for added := 0; added < days; {
		current = current.AddDate(0, 0, 1)
		if current.Weekday() != time.Saturday &&
			current.Weekday() != time.Sunday {
			added++
		}
	}
	return current
}

func (input *UpdateJobApplicationInput) IsEmpty() bool {
	return input.AppliedOn == nil &&
		input.CompanyName == nil &&
		input.RoleTitle == nil &&
		input.DatePosted == nil &&
		input.Description == nil &&
		input.JobURL == nil &&
		input.Notes == nil &&
		input.NextFollowUpOn == nil &&
		input.Location == nil &&
		input.PostedAt == nil &&
		input.Stage == nil &&
		input.User == nil &&
		input.WorkArrangement == nil
}
