package followup

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrCannotConvertToFollowUp error = errors.New("cannot convert to follow up")
var ErrFollowUpDuplicate error = errors.New("duplicate follow up")
var ErrFollowUpViolation error = errors.New("postgres error")
var ErrFollowUpFailedToInsert error = errors.New("failed to insert follow up")
var ErrFollowUpNotFound error = errors.New("follow up not found")
var ErrFollowUpNotRead error = errors.New("follow up not read")
var ErrFollowUpFailedToUpdate error = errors.New("failed to update follow up")
var ErrFollowUpFailedToDelete error = errors.New("failed to delete follow up")
var ErrFollowUpFailedToList error = errors.New("failed to list follow ups")

var ExampleUUID = uuid.New()

type FollowUpStatus string

const (
	StatusPending   = "pending"
	StatusComplete  = "complete"
	StatusCancelled = "cancelled"
)

type FollowUp struct {
	ID               uuid.UUID
	JobApplicationID uuid.UUID
	DueOn            time.Time
	FollowUpStatus   FollowUpStatus

	MessageDraft *string
	SentAt       *time.Time
	CompletedAt  *time.Time
}
