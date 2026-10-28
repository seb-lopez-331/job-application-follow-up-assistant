package jobapplication

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var ErrJobApplicationDuplicate error = errors.New("duplicate job application")
var ErrJobApplicationFailedToInsert error = errors.New("failed to insert job application")
var ErrJobApplicationInvalidInput error = errors.New("invalid job application j")
var ErrJobApplicationNotFound error = errors.New("job application not found")
var ErrJobApplicationNotRead error = errors.New("job application could not be read in database")
var ErrJobApplicationFailedToDelete error = errors.New("failed to delete job application")
var ErrJobApplicationFailedToList error = errors.New("failed to list job applications")
var ErrJobApplicationFailedToUpdate error = errors.New("failed to update job application")
var ErrCannotConvertToJobApplication error = errors.New("cannot convert object to job application")

type LocationSubmission struct {
	Place string `validate:"required,iso3166_1_alpha2|postcode_iso3166_alpha2_field|numeric"`
}

type JobApplication struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	AppliedOn   time.Time
	CompanyName string
	RoleTitle   string

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

func (j *JobApplication) IsValid() bool {
	companyName := j.CompanyName
	roleTitle := j.RoleTitle
	appliedOn := j.AppliedOn

	if strings.TrimSpace(companyName) == "" {
		return false
	}
	if strings.TrimSpace(roleTitle) == "" {
		return false
	}

	// It's impossible to apply after registering the entry
	if appliedOn.After(time.Now()) {
		return false
	}

	// Process optional fields
	location := j.Location
	if location != nil {
		submission := LocationSubmission{Place: *location}
		validate := validator.New()
		err := validate.Struct(submission)
		if err != nil {
			return false
		}
	}

	postedAt := j.PostedAt
	if postedAt != nil {
		if postedAt.After(time.Now()) {
			return false
		}
	}

	stage := j.Stage
	if stage != nil &&
		*stage != "saved" &&
		*stage != "applied" &&
		*stage != "recruiter_screen" &&
		*stage != "interviewing" &&
		*stage != "offer" &&
		*stage != "rejected" &&
		*stage != "withdrawn" {
		return false
	} else if stage == nil {
		str := "applied"
		stage = &str
	}

	workArrangement := j.WorkArrangement
	if *workArrangement != "hybrid" &&
		*workArrangement != "remote" &&
		*workArrangement != "on-site" {
		return false
	}

	return true
}

func ToJobApplication(j any) (*JobApplication, error) {
	switch app := j.(type) {
	case nil:
		return nil, ErrCannotConvertToJobApplication
	case *JobApplication:
		if app == nil {
			return nil, ErrCannotConvertToJobApplication
		}
		return app, nil
	case JobApplication:
		return &app, nil
	case CreateJobApplicationInput:
		return &JobApplication{
			AppliedOn:       app.AppliedOn,
			CompanyName:     app.CompanyName,
			RoleTitle:       app.RoleTitle,
			DatePosted:      app.DatePosted,
			Description:     app.Description,
			JobURL:          app.JobURL,
			Notes:           app.Notes,
			NextFollowUpOn:  app.NextFollowUpOn,
			Location:        app.Location,
			PostedAt:        app.PostedAt,
			Stage:           app.Stage,
			WorkArrangement: app.WorkArrangement,
		}, nil
	case UpdateJobApplicationInput:
		var companyName string
		var roleTitle string
		if app.CompanyName != nil {
			companyName = *app.CompanyName
		}
		if app.RoleTitle != nil {
			roleTitle = *app.RoleTitle
		}
		return &JobApplication{
			AppliedOn:       time.Now(),
			CompanyName:     companyName,
			RoleTitle:       roleTitle,
			DatePosted:      app.DatePosted,
			Description:     app.Description,
			JobURL:          app.JobURL,
			Notes:           app.Notes,
			NextFollowUpOn:  app.NextFollowUpOn,
			Location:        app.Location,
			PostedAt:        app.PostedAt,
			Stage:           app.Stage,
			WorkArrangement: app.WorkArrangement,
		}, nil
	case JobApplicationFilter:
		var companyName string
		var roleTitle string
		if app.CompanyName != nil {
			companyName = *app.CompanyName
		}
		if app.RoleTitle != nil {
			roleTitle = *app.RoleTitle
		}
		return &JobApplication{
			AppliedOn:       time.Now(),
			CompanyName:     companyName,
			RoleTitle:       roleTitle,
			DatePosted:      app.DatePosted,
			Description:     app.Description,
			JobURL:          app.JobURL,
			Notes:           app.Notes,
			NextFollowUpOn:  app.NextFollowUpOn,
			Location:        app.Location,
			PostedAt:        app.PostedAt,
			Stage:           app.Stage,
			WorkArrangement: app.WorkArrangement,
		}, nil
	}
	return nil, ErrCannotConvertToJobApplication
}
