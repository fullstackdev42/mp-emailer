package api

import (
	"os"
	"strconv"

	"github.com/jonesrussell/mp-emailer/campaign"
	"github.com/jonesrussell/mp-emailer/logger"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/jonesrussell/mp-emailer/user"
	"go.uber.org/fx"
)

// HandlerParams are the parameters for the API handler
type HandlerParams struct {
	fx.In

	CampaignService campaign.ServiceInterface
	UserService     user.ServiceInterface
	Logger          logger.Interface
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
		func(base *Handler, log logger.Interface) *Handler {
			return NewLoggingHandlerDecorator(base, log)
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
