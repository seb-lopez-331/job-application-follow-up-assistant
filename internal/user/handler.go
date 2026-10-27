package user

import (
	"context"

	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/constants"
)

func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(constants.UserKey).(*User)
	return u, ok
}
