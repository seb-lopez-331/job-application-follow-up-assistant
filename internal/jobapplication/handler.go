package jobapplication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/httpresponse"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/user"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service JobApplicationService
}

func NewHandler(s JobApplicationService) handler {
	return handler{service: s}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	var input CreateJobApplicationInput
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	if err := jsonDecoder.Decode(&input); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	app, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, user.ErrForbidden) {
			httpresponse.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, ErrJobApplicationDuplicate) {
			httpresponse.WriteError(w, http.StatusConflict, "job application already exists")
			return
		}
		if errors.Is(err, ErrJobApplicationInvalidInput) {
			httpresponse.WriteError(w, http.StatusUnprocessableEntity, "invalid input")
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, app)
}

func (h *handler) Read(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	app, err := h.service.Read(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrJobApplicationNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, "job application not found")
			return
		}
		if errors.Is(err, user.ErrForbidden) {
			httpresponse.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		/* internal database error */
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, app)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	var input UpdateJobApplicationInput
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	if err := jsonDecoder.Decode(&input); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	app, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrJobApplicationInvalidInput) {
			httpresponse.WriteError(w, http.StatusUnprocessableEntity, "invalid json")
			return
		}
		if errors.Is(err, ErrJobApplicationNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, "job application not found")
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, app)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrForbidden) {
			httpresponse.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	var decoder = schema.NewDecoder()
	decoder.RegisterConverter(time.Time{}, func(value string) reflect.Value {
		if t, err := time.Parse(time.DateOnly, value); err == nil {
			return reflect.ValueOf(t)
		}
		return reflect.ValueOf(time.Time{})
	})

	var filter JobApplicationFilter
	queryParams := r.URL.Query()
	err = decoder.Decode(&filter, queryParams)
	if err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("invalid query parameter: %v", err),
		)
		return
	}

	queriedUserIDStr := queryParams.Get("user_id")
	if queriedUserIDStr != "" {
		queriedUserID, err := uuid.Parse(queriedUserIDStr)
		if err != nil {
			httpresponse.WriteError(
				w,
				http.StatusBadRequest,
				fmt.Sprintf("invalid query parameter: %v", err),
			)
			return
		}

		u, _ := user.FromContext(r.Context())
		if !u.IsAdmin && u.ID != queriedUserID {
			httpresponse.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	apps, err := h.service.List(r.Context(), filter)
	if err != nil {
		// return the whole thing and ignore the bad filter
		if errors.Is(err, ErrJobApplicationInvalidInput) {
			apps, err = h.service.List(r.Context(), JobApplicationFilter{})
			if err != nil {
				httpresponse.WriteError(w, http.StatusInternalServerError, "internal error")
				return
			}
			httpresponse.WriteJSON(w, http.StatusOK, apps)
			return
		}
		httpresponse.WriteJSON(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, apps)
}

func authenticatedUserID(ctx context.Context) (*uuid.UUID, error) {
	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	return &u.ID, nil
}
