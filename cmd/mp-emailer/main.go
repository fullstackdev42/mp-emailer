package main

import (
	"fmt"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/handlers"
	"github.com/jonesrussell/loggo"
)

func main() {
	logger, err := loggo.NewLogger("mp-emailer.log", loggo.LevelInfo)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a new handler with the logger
	h := handlers.NewHandler(logger)

	// Register routes
	mux.HandleFunc("GET /", h.HandleIndex)
	mux.HandleFunc("POST /submit", h.HandleSubmit)

	logger.Info("Starting server on :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error("Error starting server", err)
	}
}
