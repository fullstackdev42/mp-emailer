package api

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// Module is the API module
//
//nolint:gochecknoglobals
var Module = fx.Module("api",
	fx.Provide(
		NewHandler,
	),
)

// HandlerParams are the parameters for the API handler
type HandlerParams struct {
	fx.In

	CampaignService campaign.ServiceInterface
	UserService     user.ServiceInterface
	Logger          loggo.LoggerInterface
	ErrorHandler    *shared.ErrorHandler
}

// NewHandler creates a new API handler
func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		campaignService: params.CampaignService,
		userService:     params.UserService,
		logger:          params.Logger,
		errorHandler:    params.ErrorHandler,
	}
}
