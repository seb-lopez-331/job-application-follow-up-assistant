package followup

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/httpresponse"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/jobapplication"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/user"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	ListDue(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service FollowUpService
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	if _, exists := user.FromContext(r.Context()); !exists {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	input := CreateFollowUpInput{JobApplicationID: id}
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	if err := jsonDecoder.Decode(&input); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	f, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, jobapplication.ErrJobApplicationNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, "job application not found")
			return
		}
		if errors.Is(err, ErrFollowUpDuplicate) {
			httpresponse.WriteError(w, http.StatusConflict, "follow up already exists")
			return
		}
		if errors.Is(err, ErrFollowUpViolation) {
			httpresponse.WriteError(w, http.StatusUnprocessableEntity, "follow up unprocessable")
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, f)
}

func (h *handler) Read(w http.ResponseWriter, r *http.Request) {
	if _, exists := user.FromContext(r.Context()); !exists {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	f, err := h.service.Read(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrForbidden) {
			httpresponse.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, ErrFollowUpNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, f)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	if _, exists := user.FromContext(r.Context()); !exists {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	input := UpdateFollowUpInput{}
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	if err := jsonDecoder.Decode(&input); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	f, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, user.ErrForbidden) {
			httpresponse.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, ErrFollowUpNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal error")
	}

	httpresponse.WriteJSON(w, http.StatusOK, f)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	if _, exists := user.FromContext(r.Context()); !exists {
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
	if err != nil && !errors.Is(err, ErrFollowUpNotFound) {
		if errors.Is(err, user.ErrForbidden) {
			httpresponse.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) ListDue(w http.ResponseWriter, r *http.Request) {
	if _, exists := user.FromContext(r.Context()); !exists {
		httpresponse.WriteError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	params := r.URL.Query()
	queriedUserIDStr := params.Get("user_id")
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

	followUps, err := h.service.ListDue(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, followUps)
}
