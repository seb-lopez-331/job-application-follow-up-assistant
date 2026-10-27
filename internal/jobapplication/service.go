package jobapplication

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidInput = errors.New("invalid job application input")

type CreateJobApplicationInput struct {
	CompanyName string
	RoleTitle	string
	JobURL      string
	Status		string
	AppliedOn	*time.Time
}

// The service depends on the Repository interface, not directly on PostgreSQL.
// That makes it possible to unit-test scheduling and validation rules with
// a fake repository
type Service interface {
	Create(ctx context.Context, userID uuid.UUID,
		input CreateJobApplicationInput) (JobApplication, error)
	
	Get(ctx context.Context, userID, id uuid.UUID) (JobApplication, error)

	Update(ctx context.Context, userID, id uuid.UUID,
		input UpdateJobApplicationInput) (JobApplication, error)

	List(ctx context.Context, userID uuid.UUID,
		filter JobApplicationFilter) ([]JobApplication, error)
}

type JobApplicationService struct {
	repo Repository
}

func NewJobApplicationService(repo Repository) *JobApplicationService {
	return &JobApplicationService{ repo }
}

func (s *JobApplicationService) Create(
	ctx context.Context,
	userID uuid.UUID,
	input CreateJobApplicationInput,
) (JobApplication, error) {
	if strings.TrimSpace(input.CompanyName) == "" ||
		strings.TrimSpace(input.RoleTitle) == "" {
			return JobApplication{}, ErrInvalidInput
	}

	stage := input.Stage
	if status == "" {
		status = "applied"
	}

	if stage != "saved" && stage != "applied" && stage != "recruiter_screen" 
		&& stage != "interviewing" && stage != "offer" && stage != "rejected" 
		&& stage != "withdrawn" {
		return JobApplication{}, ErrInvalidInput
	}

	app := JobApplication{
		UserID: 		userID,
		CompanyName: 	strings.TrimSpace(input.CompanyName),
		RoleTitle: 		strings.TrimSpace(input.RoleTitle),
		Stage: 			stage,
	}

	if stage == "applied" && input.AppliedOn != nil {
		due := addBusinessDays(*input.AppliedOn, 7)
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
