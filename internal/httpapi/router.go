package httpapi

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/httpresponse"
	"github.com/seb-lopez-331/job-application-follow-up-assistant/internal/jobapplication"
)

func RegisterRoutes(mux *chi.Mux, handler jobapplication.Handler) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(func(authRouter chi.Router) {
		authRouter.Use(SessionAuthMiddleware)

		authRouter.Get("/api/v1/users/me"), func(w http.ResponseWriter, r *http.Request) {
			session, ok := r.Context().Value(constants.SessionUserKey)
			if !ok {
			}
		})
	})

	// job application related routes
	r.Post("/api/v1/applications", handler.Create)

}

func SessionAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				httpresponse.WriteError(
					w,
					http.StatusUnauthorized,
					"Unauthorized: No session token provided",
				)
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}

		sessionToken := cookie.Value

	})
}
