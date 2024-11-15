package campaign

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

type MPLookupServiceParams struct {
	fx.In

	Logger loggo.LoggerInterface
}

type MPLookupService struct {
	logger loggo.LoggerInterface
}

func NewMPLookupService(params MPLookupServiceParams) *MPLookupService {
	return &MPLookupService{
		logger: params.Logger,
	}
}

func (s *MPLookupService) FetchRepresentatives(postalCode string) ([]Representative, error) {
	url := fmt.Sprintf("https://represent.opennorth.ca/postcodes/%s/?format=json", postalCode)
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
