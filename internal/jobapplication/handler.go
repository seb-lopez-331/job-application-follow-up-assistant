package jobapplication

import (
	"json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input CreateJobApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	userID := authenticatedUserID(r.Context())

	app, err := h.service.Create(
		r.Context(),
		userID,
		toCreateInput(input),
	)
	if err != nil {
		if errors.Is(err, application.ErrInvalidInput) {
			writeError(w, http.StatusUnprocessableEntity, "invalid_input")
			return
		}

		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusCreated, app)
}

func authenticatedUserID(ctx context.Context) uuid.UUID {

}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	
}