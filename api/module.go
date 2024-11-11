package api

import (
	"os"
	"strconv"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// HandlerParams are the parameters for the API handler
type HandlerParams struct {
	fx.In

	CampaignService campaign.ServiceInterface
	UserService     user.ServiceInterface
	Logger          loggo.LoggerInterface
	ErrorHandler    shared.ErrorHandlerInterface
	JWTExpiry       int
}

// Module is the API module
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		NewHandler,
		provideJWTExpiry,
	),
	fx.Decorate(
		func(base *Handler, logger loggo.LoggerInterface) *Handler {
			return NewLoggingHandlerDecorator(base, logger)
		},
	),
)

// NewHandler creates a new API handler
func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		campaignService: params.CampaignService,
		userService:     params.UserService,
		logger:          params.Logger,
		errorHandler:    params.ErrorHandler,
		jwtExpiry:       params.JWTExpiry,
	}
}

func provideJWTExpiry() (int, error) {
	expiryStr := os.Getenv("JWT_EXPIRY")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil {
		return 15, nil // Default to 15 minutes if not set or invalid
	}
	return expiry, nil
}
