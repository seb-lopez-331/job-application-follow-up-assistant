package user

import (
	"context"
	"net/http"

	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/constants"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service UserService
}

func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(constants.UserKey).(*User)
	return u, ok
}

func ValidateAdmin(ctx context.Context) *User {
	u, ok := FromContext(ctx)
	if !ok {
		return nil
	}
	return u
}
