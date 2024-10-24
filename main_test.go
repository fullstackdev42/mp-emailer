package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/assert"
)

func TestStaticFileServing(t *testing.T) {
	// Setup
	cfg := &config.Config{}
	logger, _ := loggo.NewLogger("test.log", loggo.LevelDebug)
	tmplManager, _ := server.NewTemplateManager(templateFS)
	e := newEcho(cfg, logger, tmplManager)

	// Test
	req := httptest.NewRequest(http.MethodGet, "/static/css/styles.css", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/css")
}
