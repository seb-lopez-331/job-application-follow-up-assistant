package jobapplication

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/user"
)

var ErrInvalidInput = errors.New("invalid job application input")

type LocationSubmission struct {
	Place string `validate:"required,iso3166_1_alpha2|postcode_iso3166_alpha2_field|numeric"`
}

type Job struct {
	ID          uuid.UUID
	CompanyName string
	RoleTitle   string
}

type JobApplication struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Job            Job
	Stage          string
	NextFollowUpOn *time.Time
}

type CreateJobApplicationInput struct {
	CompanyName string
	RoleTitle   string
	AppliedOn   time.Time

	Description     *string
	Notes           *string
	JobURL          *url.URL
	Location        *string
	PostedAt        *time.Time
	Stage           *string
	User            *user.User
	WorkArrangement string
}
