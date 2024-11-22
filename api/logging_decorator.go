package api

import "github.com/jonesrussell/mp-emailer/logger"

func NewLoggingHandlerDecorator(base *Handler, logger logger.Interface) *Handler {
	// Add logging decoration logic here
	return &Handler{
		campaignService: base.campaignService,
		userService:     base.userService,
		logger:          logger,
		errorHandler:    base.errorHandler,
		jwtExpiry:       base.jwtExpiry,
	}
}
