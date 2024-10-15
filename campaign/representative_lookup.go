package campaign

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jonesrussell/loggo"
)

// RepresentativeLookupServiceInterface defines the interface for representative lookup
type RepresentativeLookupServiceInterface interface {
	FetchRepresentatives(postalCode string) ([]Representative, error)
	FilterRepresentatives(representatives []Representative, filters map[string]string) []Representative
}

// RepresentativeLookupService implements RepresentativeLookupServiceInterface
type RepresentativeLookupService struct {
	logger  loggo.LoggerInterface
	baseURL string
}

func NewRepresentativeLookupService(logger loggo.LoggerInterface) RepresentativeLookupServiceInterface {
	return &RepresentativeLookupService{
		logger:  logger,
		baseURL: "https://represent.opennorth.ca",
	}
}

func (s *RepresentativeLookupService) FetchRepresentatives(postalCode string) ([]Representative, error) {
	url := fmt.Sprintf("%s/postcodes/%s/?format=json", s.baseURL, postalCode)
	s.logger.Info("Making request to", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error making request", err)
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	s.logger.Info("Response received", "status", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Error reading response body", err)
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		s.logger.Error("Error unmarshaling JSON", err)
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	return apiResp.RepresentativesCentroid, nil
}

func (s *RepresentativeLookupService) FilterRepresentatives(representatives []Representative, filters map[string]string) []Representative {
	filtered := make([]Representative, 0)
	for _, rep := range representatives {
		match := true
		for key, value := range filters {
			switch key {
			case "type":
				if rep.ElectedOffice != value {
					match = false
				}
			case "level":
				if rep.ElectedOffice != value {
					match = false
				}
			case "name":
				if !strings.Contains(strings.ToLower(rep.Name), strings.ToLower(value)) {
					match = false
				}
			case "party":
				if !strings.EqualFold(rep.Party, value) {
					match = false
				}
			default:
				s.logger.Warn("Unknown filter key", "key", key)
			}
		}
		if match {
			filtered = append(filtered, rep)
		}
	}
	return filtered
}
