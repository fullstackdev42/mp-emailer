package campaign

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jonesrussell/mp-emailer/logger"
	"go.uber.org/fx"
)

type MPLookupServiceParams struct {
	fx.In

	Logger logger.Interface
}

type MPLookupService struct {
	Logger logger.Interface
}

func NewMPLookupService(params MPLookupServiceParams) *MPLookupService {
	return &MPLookupService{
		Logger: params.Logger,
	}
}

func (s *MPLookupService) FetchRepresentatives(postalCode string) ([]Representative, error) {
	url := fmt.Sprintf("https://represent.opennorth.ca/postcodes/%s/?format=json", postalCode)
	s.Logger.Info("Making request to", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		s.Logger.Error("Error making request", err)
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.Logger.Error("Error closing response body", err)
		}
	}(resp.Body)

	s.Logger.Info("Response received", "status", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error("Error reading response body", err)
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		s.Logger.Error("Error unmarshaling JSON", err)
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return apiResp.RepresentativesCentroid, nil
}
