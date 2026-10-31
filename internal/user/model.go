package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rickar/cal/v2"
)

var ErrUserNotFoundInCtx = errors.New("user not found in context")
var ErrForbidden = errors.New("forbidden")

type FollowUpSchedule struct {
	FollowUpDaysAfter []int
	SkipWeekends      bool
	SkipHolidays      bool
	Auto              bool
	SelectedHolidays  *cal.Calendar
}

type Preferences struct {
	FollowUpSchedule FollowUpSchedule
	Mode             string // light, dark, or system
	ItemsPerPage     int
	Notifications    bool
}

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	IsAdmin      bool
	Preferences  Preferences
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
