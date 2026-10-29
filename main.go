package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hotels-api/config"
	"hotels-api/models"
	"hotels-api/routes"
	"hotels-api/seeder"
)

func main() {
	log.Println("[INFO] Initializing Hotels Multi-Tenant OTA Backend API...")

	// 1. Connect Database
	db := config.ConnectDB()

	// 2. Auto-migrate schema
	if err := models.MigrateDB(db); err != nil {
		log.Fatalf("[FATAL] Database auto-migration failed: %v", err)
	}
	log.Println("[INFO] Database schema migration verified successfully")

	// 3. Seed initial dataset
	seeder.Seed()

	// 4. Setup Routes
	router := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 5. Graceful shutdown routine
	go func() {
		log.Printf("[INFO] Hotels Backend API server listening on http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] HTTP server failed to listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[FATAL] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Hotels Backend API exited cleanly")
}
