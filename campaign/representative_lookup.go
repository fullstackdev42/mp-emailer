package campaign

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
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

type RepresentativeLookupServiceParams struct {
	fx.In

	Config *config.Config
	Logger loggo.LoggerInterface
}

// NewRepresentativeLookupService creates a new instance of RepresentativeLookupService
func NewRepresentativeLookupService(params RepresentativeLookupServiceParams) RepresentativeLookupServiceInterface {
	return &RepresentativeLookupService{
		logger:  params.Logger,
		baseURL: params.Config.Server.RepresentativeLookupBaseURL,
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error("Error closing response body", err)
		}
	}(resp.Body)

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
		if s.matchesFilters(rep, filters) {
			filtered = append(filtered, rep)
		}
	}
	return filtered
}

func (s *RepresentativeLookupService) matchesFilters(rep Representative, filters map[string]string) bool {
	for key, value := range filters {
		switch key {
		case "type":
			if rep.ElectedOffice != value {
				return false
			}
		case "level":
			if rep.ElectedOffice != value {
				return false
			}
		case "name":
			if !strings.Contains(strings.ToLower(rep.Name), strings.ToLower(value)) {
				return false
			}
		case "party":
			if !strings.EqualFold(rep.Party, value) {
				return false
			}
		default:
			s.logger.Warn("Unknown filter key", "key", key)
		}
	}
	return true
}
