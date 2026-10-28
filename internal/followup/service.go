package followup

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/jobapplication"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/user"
)

var ErrInvalidInput = errors.New("invalid followup input")
var ErrCannotGetApplicationFromFollowUp = errors.New("cannot get job application from follow up")

type Service interface {
	Create(ctx context.Context, input CreateFollowUpInput) (*FollowUp, error)
	Read(ctx context.Context, id uuid.UUID) (*FollowUp, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateFollowUpInput) (*FollowUp, error)
	Delete(ctx context.Context, id uuid.UUID) error

	Schedule(ctx context.Context, appID uuid.UUID, dueOn time.Time) (*FollowUp, error)
	Complete(ctx context.Context, followUpID uuid.UUID) error
	ListDue(ctx context.Context, dueOn time.Time) ([]FollowUp, error)
}

type FollowUpService struct {
	repo                  Repository
	jobApplicationService jobapplication.JobApplicationService
}

type CreateFollowUpInput struct {
	DueOn            time.Time
	JobApplicationID uuid.UUID

	MessageDraft *string
}

type UpdateFollowUpInput struct {
	DueOn          *time.Time
	MessageDraft   *string
	FollowUpStatus *string
	SentAt         *time.Time
	CompletedAt    *time.Time
}

func (s *FollowUpService) Create(
	ctx context.Context,
	input CreateFollowUpInput,
) (*FollowUp, error) {
	// Verify whether the application id passed in is valid
	if _, err := s.jobApplicationService.Read(ctx, input.JobApplicationID); err != nil {
		return nil, err
	}

	now := time.Now()
	if input.DueOn.Before(now) {
		return nil, ErrInvalidInput
	}

	return s.repo.Create(ctx, input)
}

func (s *FollowUpService) Read(ctx context.Context, id uuid.UUID) (*FollowUp, error) {
	f, err := s.repo.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	app, err := s.jobApplicationService.Read(ctx, f.JobApplicationID)
	if err != nil { // follow up exists but job application is inaccessible, return 500
		return nil, ErrCannotGetApplicationFromFollowUp
	}

	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}

	if !u.IsAdmin && u.ID != app.UserID {
		return nil, user.ErrForbidden
	}

	return f, nil
}

func (s *FollowUpService) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateFollowUpInput,
) (*FollowUp, error) {
	// Check if the user has access
	_, err := s.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.IsEmpty() {
		return nil, ErrInvalidInput
	}

	// Validate input
	// A follow up can't be updated to be due before the time the request was made
	now := time.Now()
	if input.DueOn != nil && (*input.DueOn).Before(now) {
		return nil, ErrInvalidInput
	}
	// A follow up can't be updated to be sent at after the time the request was made
	if input.SentAt != nil && (*input.SentAt).After(now) {
		return nil, ErrInvalidInput
	}
	if input.CompletedAt != nil && (*input.CompletedAt).After(now) {
		return nil, ErrInvalidInput
	}

	if input.FollowUpStatus != nil {
		if *input.FollowUpStatus != StatusPending && *input.FollowUpStatus != StatusComplete &&
			*input.FollowUpStatus != StatusCancelled {
			return nil, ErrInvalidInput
		}
		if *input.FollowUpStatus == StatusComplete && input.CompletedAt == nil {
			return nil, ErrInvalidInput
		}
	}

	return s.repo.Update(ctx, id, input)
}

func (s *FollowUpService) Schedule(
	ctx context.Context,
	appID uuid.UUID,
	dueOn time.Time,
) (*FollowUp, error) {
	input := CreateFollowUpInput{DueOn: dueOn, JobApplicationID: appID}
	return s.Create(ctx, input)
}

func (s *FollowUpService) Complete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	complete := StatusComplete
	_, err := s.Update(ctx, id, UpdateFollowUpInput{
		CompletedAt:    &now,
		FollowUpStatus: &complete,
	})
	return err
}

func (s *FollowUpService) ListDue(ctx context.Context, dueOn time.Time) ([]FollowUp, error) {
	return s.ListDue(ctx, dueOn)
}

func (input *UpdateFollowUpInput) IsEmpty() bool {
	return input.DueOn == nil &&
		input.MessageDraft == nil &&
		input.FollowUpStatus == nil &&
		input.SentAt == nil &&
		input.CompletedAt == nil
}
