package api

import "github.com/jonesrussell/loggo"

func NewLoggingHandlerDecorator(base *Handler, logger loggo.LoggerInterface) *Handler {
	// Add logging decoration logic here
	return &Handler{
		campaignService: base.campaignService,
		userService:     base.userService,
		logger:          logger,
		errorHandler:    base.errorHandler,
		jwtExpiry:       base.jwtExpiry,
	}
}
