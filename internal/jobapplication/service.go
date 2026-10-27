package jobapplication

import (
	"context"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/user"
)

// The service depends on the Repository interface, not directly on PostgreSQL.
// That makes it possible to unit-test scheduling and validation rules with
// a fake repository
type JobService interface {
	Create(ctx context.Context, userID uuid.UUID,
		input CreateJobInput) (*Job, error)

	Get(ctx context.Context, id uuid.UUID) (*Job, error)

	GetByCompanyNameAndJobTitle(ctx context.Context,
		companyName, jobTitle string) (*Job, error)

	Update(ctx context.Context, id uuid.UUID,
		input UpdateJobInput) (*Job, error)

	List(ctx context.Context, userID uuid.UUID,
		filter JobFilter) ([]Job, error)
}

type JobApplicationService interface {
	Create(ctx context.Context, userID uuid.UUID,
		input CreateJobApplicationInput) (*JobApplication, error)

	Get(ctx context.Context, userID, id uuid.UUID) (*JobApplication, error)

	UpdateJobApplication(ctx context.Context, userID, id uuid.UUID,
		input UpdateJobApplicationInput) (*JobApplication, error)

	List(ctx context.Context, userID uuid.UUID,
		filter JobApplicationFilter) ([]JobApplication, error)
}

type jobService struct {
	repo Repository
}

type jobApplicationService struct {
	jobService jobService
	repo       Repository
}

func NewJobApplicationService(
	repo Repository,
	jobService jobService,
) *jobApplicationService {
	return &jobApplicationService{jobService, repo}
}

func (s *jobService) GetByCompanyNameAndJobTitle(
	ctx context.Context,
	companyName, jobTitle string,
) (*Job, error) {
	if companyName == "" || jobTitle == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.GetByCompanyNameAndJobTitle(ctx, companyName, jobTitle)
}

func (s *jobApplicationService) Create(
	ctx context.Context,
	userID uuid.UUID,
	input CreateJobApplicationInput,
) (*JobApplication, error) {
	if strings.TrimSpace(input.CompanyName) == "" ||
		strings.TrimSpace(input.RoleTitle) == "" {
		return nil, ErrInvalidInput
	}

	companyName := input.CompanyName
	roleTitle := input.RoleTitle
	appliedOn := input.AppliedOn

	if strings.TrimSpace(companyName) == "" {
		return nil, ErrInvalidInput
	}
	if strings.TrimSpace(roleTitle) == "" {
		return nil, ErrInvalidInput
	}

	// It's impossible to apply after registering the entry
	if appliedOn.After(time.Now()) {
		return nil, ErrInvalidInput
	}

	// Process optional fields
	location := input.Location
	if location != nil {
		submission := LocationSubmission{Place: *location}
		validate := validator.New()
		err := validate.Struct(submission)
		if err != nil {
			return nil, ErrInvalidInput
		}
	}

	postedAt := input.PostedAt
	if postedAt != nil {
		if postedAt.After(time.Now()) {
			return nil, ErrInvalidInput
		}
	}

	stage := input.Stage
	if stage != nil &&
		*stage != "saved" &&
		*stage != "applied" &&
		*stage != "recruiter_screen" &&
		*stage != "interviewing" &&
		*stage != "offer" &&
		*stage != "rejected" &&
		*stage != "withdrawn" {
		return nil, ErrInvalidInput
	} else if stage == nil {
		str := "applied"
		stage = &str
	}

	if input.User != nil {
		userFromContext, ok := user.UserFromContext(ctx)
		if !ok {
			return nil, user.ErrUserNotFoundInCtx
		}
		if !userFromContext.IsAdmin {
			return nil, user.ErrForbidden
		}
	}

	workArrangement := input.WorkArrangement
	if workArrangement != "hybrid" &&
		workArrangement != "remote" &&
		workArrangement != "on-site" {
		return nil, ErrInvalidInput
	}

	// Get existing job in the DB
	job, err := s.jobService.GetByCompanyNameAndJobTitle(ctx, companyName, roleTitle)

	if *stage == "applied" {
		due := addBusinessDays(input.AppliedOn, 7)
		app.NextFollowUpOn = &due
	}

	return s.repo.Create(ctx, app)
}

func (s *JobApplicationService) Get(
	ctx context.Context,
	userID, id uuid.UUID,
) (JobApplication, error) {

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
