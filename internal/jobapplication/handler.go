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

var decoder = schema.NewDecoder()

func init() {
	// Register a universal layout handler for time.Time (e.g., "2026-10-27")
	decoder.RegisterConverter(time.Time{}, func(value string) reflect.Value {
		if t, err := time.Parse(time.DateOnly, value); err == nil {
			return reflect.ValueOf(t)
		}
		return reflect.ValueOf(time.Time{})
	})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	var input CreateJobApplicationInput
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	if err := jsonDecoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	app, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, user.ErrForbidden) {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, ErrJobApplicationDuplicate) {
			writeError(w, http.StatusConflict, "job application already exists")
			return
		}
		if errors.Is(err, ErrJobApplicationInvalidInput) {
			writeError(w, http.StatusUnprocessableEntity, "invalid input")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, app)
}

func (h *handler) Read(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	app, err := h.service.Read(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrJobApplicationNotFound) {
			writeError(w, http.StatusNotFound, "job application not found")
			return
		}
		if errors.Is(err, user.ErrForbidden) {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}
		/* internal database error */
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, app)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	var input UpdateJobApplicationInput
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	if err := jsonDecoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	app, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrJobApplicationInvalidInput) {
			writeError(w, http.StatusUnprocessableEntity, "invalid json")
			return
		}
		if errors.Is(err, ErrJobApplicationNotFound) {
			writeError(w, http.StatusNotFound, "job application not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, app)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrForbidden) {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	_, err := authenticatedUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	var filter JobApplicationFilter
	queryParams := r.URL.Query()
	err = decoder.Decode(&filter, queryParams)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid query parameter: %v", err))
		return
	}

	queriedUserIDStr := queryParams.Get("user_id")
	if queriedUserIDStr != "" {
		queriedUserID, err := uuid.Parse(queriedUserIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid query parameter: %v", err))
			return
		}

		u, _ := user.FromContext(r.Context())
		if !u.IsAdmin && u.ID != queriedUserID {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	apps, err := h.service.List(r.Context(), filter)
	if err != nil {
		// return the whole thing and ignore the bad filter
		if errors.Is(err, ErrJobApplicationInvalidInput) {
			apps, err = h.service.List(r.Context(), JobApplicationFilter{})
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal error")
				return
			}
			writeJSON(w, http.StatusOK, apps)
			return
		}
		writeJSON(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, apps)
}

func authenticatedUserID(ctx context.Context) (*uuid.UUID, error) {
	u, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotFoundInCtx
	}
	return &u.ID, nil
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(struct {
		Data any `json:"data"`
	}{Data: v})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
