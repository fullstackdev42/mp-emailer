package user

import (
	"fmt"

	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// getLogger retrieves the logger from the context
func getLogger(c echo.Context) (loggo.LoggerInterface, error) {
	logger, ok := c.Get("logger").(loggo.LoggerInterface)
	if !ok {
		return nil, fmt.Errorf("logger not found in context")
	}
	return logger, nil
}

// GetOwnerIDFromSession retrieves the owner ID from the session
func GetOwnerIDFromSession(c echo.Context) (string, error) {
	logger, err := getLogger(c)
	if err != nil {
		return "", err
	}
	logger.Debug("GetOwnerIDFromSession: Starting")

	ownerID, ok := c.Get("user_id").(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in session or not a string")
	}

	logger.Debug("GetOwnerIDFromSession: Owner ID retrieved", "ownerID", ownerID)
	return ownerID, nil
}
