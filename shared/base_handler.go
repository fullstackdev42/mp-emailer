package shared

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// BaseHandlerParams defines the input parameters for BaseHandler
type BaseHandlerParams struct {
	fx.In

	Store        sessions.Store
	Logger       loggo.LoggerInterface
	ErrorHandler ErrorHandlerInterface
	Config       *config.Config
}

type BaseHandler struct {
	Store        sessions.Store
	Logger       loggo.LoggerInterface
	ErrorHandler ErrorHandlerInterface
	Config       *config.Config
}

func NewBaseHandler(params BaseHandlerParams) BaseHandler {
	return BaseHandler{
		Store:        params.Store,
		Logger:       params.Logger,
		ErrorHandler: params.ErrorHandler,
		Config:       params.Config,
	}
}
