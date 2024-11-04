package campaign

import (
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

type Handler struct {
	shared.BaseHandler
	service                     ServiceInterface
	representativeLookupService RepresentativeLookupServiceInterface
	emailService                email.Service
	client                      ClientInterface
	mapError                    func(error) (int, string)
}

// HandlerParams for dependency injection
type HandlerParams struct {
	shared.BaseHandlerParams
	fx.In
	Service                     ServiceInterface
	Logger                      loggo.LoggerInterface
	RepresentativeLookupService RepresentativeLookupServiceInterface
	EmailService                email.Service
	Client                      ClientInterface
	ErrorHandler                *shared.ErrorHandler
	TemplateRenderer            shared.TemplateRendererInterface
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler initializes a new Handler
func NewHandler(params HandlerParams) (HandlerResult, error) {
	handler := &Handler{
		BaseHandler:                 shared.NewBaseHandler(params.BaseHandlerParams),
		service:                     params.Service,
		representativeLookupService: params.RepresentativeLookupService,
		emailService:                params.EmailService,
		client:                      params.Client,
		mapError:                    mapErrorToHTTPStatus,
	}
	return HandlerResult{Handler: handler}, nil
}
