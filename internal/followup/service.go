package followup

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidInput = errors.New("invalid followup input")

type Service interface {
	Schedule(ctx context.Context, userID, applicationID uuid.UUID,
		dueOn time.Time) (FollowUp, error)
	
	Complete(ctx context.Context, userID, followUpID uuid.UUID) error

	ListDue(ctx context.Context, userID uuid.UUID,
		dueOn time.Time) ([]FollowUp, error)
}

type FollowUpService struct {
	repo Repository
}
