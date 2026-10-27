package followup

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FollowUp struct {
	ID					uuid.UUID
	JobApplicationID	uuid.UUID
	DueOn				*time.Time
	MessageDraft		string
	FollowUpStatus		string
	SentAt				*time.Time
	CompletedAt			*time.Time
}
