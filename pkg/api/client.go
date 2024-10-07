package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/jonesrussell/loggo"
)

type Client struct {
	logger loggo.LoggerInterface
}

func NewClient(logger loggo.LoggerInterface) *Client {
	return &Client{logger: logger}
}

func (c *Client) FetchRepresentatives(postalCode string) ([]models.Representative, error) {
	url := fmt.Sprintf("https://represent.opennorth.ca/postcodes/%s/?format=json", postalCode)
	c.logger.Info("Making request to", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		c.logger.Error("Error making request", err)
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	c.logger.Info("Response received", "status", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Error reading response body", err)
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var apiResp models.APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		c.logger.Error("Error unmarshaling JSON", err)
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return apiResp.RepresentativesCentroid, nil
}
