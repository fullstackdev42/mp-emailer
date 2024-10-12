package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Logger          loggo.LoggerInterface
	client          api.ClientInterface
	Store           sessions.Store
	emailService    services.EmailService
	templateManager *templates.TemplateManager
}

func NewHandler(logger loggo.LoggerInterface, client api.ClientInterface, store sessions.Store, emailService services.EmailService, tmplManager *templates.TemplateManager) *Handler {

	return &Handler{
		Logger:          logger,
		client:          client,
		Store:           store,
		emailService:    emailService,
		templateManager: tmplManager,
	}
}

func (h *Handler) HandleIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func (h *Handler) HandleSubmit(c echo.Context) error {
	h.Logger.Info("Handling submit request")

	postalCode := c.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.Logger.Warn("Invalid postal code submitted", "postalCode", postalCode)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid postal code format")
	}

	mpFinder := services.NewMPFinder(h.client, h.Logger)

	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error finding MP")
	}

	emailContent := composeEmail(mp)

	data := struct {
		Email   string
		Content string
	}{
		Email:   mp.Email,
		Content: emailContent,
	}

	return c.Render(http.StatusOK, "email.html", data)
}

func (h *Handler) HandleEcho(c echo.Context) error {
	type EchoRequest struct {
		Message string `json:"message"`
	}

	req := new(EchoRequest)
	if err := c.Bind(req); err != nil {
		return h.handleError(err, http.StatusBadRequest, "Error binding request")
	}

	return c.JSON(http.StatusOK, req)
}

func composeEmail(mp models.Representative) string {
	return fmt.Sprintf("Dear %s,\n\nThis is a sample email content.\n\nBest regards,\nYour constituent", mp.Name)
}
