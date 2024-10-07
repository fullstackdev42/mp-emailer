package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define a variable to hold the findMP function
var findMPFunc = findMP

func TestHandleIndex(t *testing.T) {
	mockLogger := new(loggo.MockLogger)
	logger = mockLogger

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleIndex)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "MP Emailer")
	assert.Contains(t, rr.Body.String(), "<form action=\"/submit\" method=\"post\">")
}

func TestHandleSubmit(t *testing.T) {
	mockLogger := new(loggo.MockLogger)
	logger = mockLogger

	mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	form := url.Values{}
	form.Add("postalCode", "K1A0A6")
	req, err := http.NewRequest("POST", "/submit", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSubmit)

	// Test case 1: MP found
	oldFindMP := findMPFunc
	defer func() { findMPFunc = oldFindMP }()
	findMPFunc = func(postalCode string) (Representative, error) {
		return Representative{
			Name:  "Test MP",
			Email: "test@example.com",
		}, nil
	}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Email to MP")
	assert.Contains(t, rr.Body.String(), "test@example.com")
	assert.Contains(t, rr.Body.String(), "Dear Test MP,")

	// Test case 2: MP not found
	findMPFunc = func(postalCode string) (Representative, error) {
		return Representative{}, fmt.Errorf("no MP found for postal code %s", postalCode)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error finding MP")

	mockLogger.AssertExpectations(t)
}

func TestFindMP(t *testing.T) {
	mockLogger := new(loggo.MockLogger)
	logger = mockLogger

	mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := PostalCodeResponse{
			Representatives: []Representative{
				{
					Name:          "Test MP",
					ElectedOffice: "MP",
					Email:         "test@example.com",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Replace the URL in findMP with the mock server URL
	oldFindMP := findMPFunc
	defer func() { findMPFunc = oldFindMP }()
	findMPFunc = func(postalCode string) (Representative, error) {
		url := server.URL + "/postcodes/" + postalCode + "/?format=json"
		resp, err := http.Get(url)
		if err != nil {
			return Representative{}, err
		}
		defer resp.Body.Close()

		var postalCodeResp PostalCodeResponse
		err = json.NewDecoder(resp.Body).Decode(&postalCodeResp)
		if err != nil {
			return Representative{}, err
		}

		for _, rep := range postalCodeResp.Representatives {
			if rep.ElectedOffice == "MP" {
				logger.Info("MP found", "name", rep.Name, "email", rep.Email)
				return rep, nil
			}
		}

		return Representative{}, nil
	}

	mp, err := findMPFunc("K1A0A6")
	assert.NoError(t, err)
	assert.Equal(t, "Test MP", mp.Name)
	assert.Equal(t, "test@example.com", mp.Email)

	mockLogger.AssertExpectations(t)
}

func TestComposeEmail(t *testing.T) {
	mp := Representative{
		Name:  "Test MP",
		Email: "test@example.com",
	}

	email := composeEmail(mp)

	assert.Contains(t, email, "Dear Test MP,")
	assert.Contains(t, email, "I am writing to express my concerns about [ISSUE].")
	assert.Contains(t, email, "Sincerely,")
}
