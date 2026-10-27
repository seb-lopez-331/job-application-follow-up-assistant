package followup

import (
	"json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

