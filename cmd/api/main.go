package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"	

	"job-application-follow-up-assistant/internal/jobapplication"
	"job-application-follow-up-assistant/internal/httpapi"
)

func main() {
	// Initalize the root context to control application lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the new database connection pool
	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("initialize database pool: %v", err)
	}
	defer db.Close()

	// Verify connectivity during startup
	pingCtx, pingCancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer pingCancel()

	if err := db.Ping(pingCtx); err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	// Wire dependencies
	repo := jobapplication.NewPostgresRepository(db)
	service := jobapplication.NewService(repo)
	handler := application.NewHandler(service)

	// Build the HTTP router
	router := chi.NewRouter()
	httpapi.RegisterRoutes(router, handler)

	server := &http.Server{
		Addr:				":8080",
		Handler:			router,
		ReadHeaderTimeout: 	5 * time.Second,
	}

	log.Printf("API listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil &&
		err !- http.ErrServerClosed {
		log.Fatalf("run HTTP server: %v", err)
	}
}