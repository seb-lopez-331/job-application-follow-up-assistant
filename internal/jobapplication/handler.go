package jobapplication

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/user"
)

type Handler struct {
	service Service
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateJobApplicationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	userID, err := authenticatedUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	app, err := h.service.Create(
		r.Context(),
		*userID,
		input,
	)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusUnprocessableEntity, "invalid_input")
			return
		}

		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusCreated, app)
}

func authenticatedUserID(ctx context.Context) (*uuid.UUID, error) {
	u, ok := user.UserFromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	return &u.ID, nil
}

func toCreateInput(
	ctx context.Context,
	input CreateJobApplicationRequest,
) (*CreateJobApplicationInput, error) {
	

	return &CreateJobApplicationInput{
		CompanyName: companyName,
		RoleTitle: roleTitle,
		JobURL: 
	}, nil

}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {

}
