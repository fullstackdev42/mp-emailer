package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func NewJSONRequest(method, path string, body interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req := httptest.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req, nil
}

func ExtractJSONResponse(rec *httptest.ResponseRecorder, v interface{}) error {
	return json.NewDecoder(rec.Body).Decode(v)
}
