package main

import (
	"fmt"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/handlers"
	"github.com/jonesrussell/loggo"
)

var logger loggo.LoggerInterface

func main() {
	var err error
	logger, err = loggo.NewLogger("mp-emailer.log", loggo.LevelInfo)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleIndex)
	mux.HandleFunc("/submit", handlers.HandleSubmit)

	logger.Info("Starting server on :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error("Error starting server", err)
	}
}
