package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/routes"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/server"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	databaseURL := "postgres://local:local@localhost:5432/fb_ingest"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Create tables if they don't exist
	if err := server.CreateTables(ctx, dbpool); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Initialize the server with routes and database pool
	srv := server.NewServer(dbpool)

	// Set up routes
	routes.SetupRoutes(srv)

	// Start the HTTP server
	port := "8080"
	fmt.Printf("Server is running on :%s\n", port)
	if err := http.ListenAndServe(":"+port, srv.Router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
