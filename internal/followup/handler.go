package followup

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
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
		writeError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	input := CreateFollowUpInput{JobApplicationID: id}
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	if err := jsonDecoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	f, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, jobapplication.ErrJobApplicationNotFound) {
			writeError(w, http.StatusNotFound, "job application not found")
			return
		}
		if errors.Is(err, ErrFollowUpDuplicate) {
			writeError(w, http.StatusConflict, "follow up already exists")
			return
		}
		if errors.Is(err, ErrFollowUpViolation) {
			writeError(w, http.StatusUnprocessableEntity, "follow up unprocessable")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (h *handler) Read(w http.ResponseWriter, r *http.Request) {
	if _, exists := user.FromContext(r.Context()); !exists {
		writeError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) ListDue(w http.ResponseWriter, r *http.Request) {

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
