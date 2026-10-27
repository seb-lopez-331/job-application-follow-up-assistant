package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFoundInCtx = errors.New("user not found in context")
var ErrForbidden = errors.New("forbidden")

type FollowUpSchedule struct {
	FollowUpDaysAfter []int
	SkipWeekends      bool
	SkipHolidays      bool
}

type Preferences struct {
	FollowUpSchedule FollowUpSchedule
	Mode             string // light or dark
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
