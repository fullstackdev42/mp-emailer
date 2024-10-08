package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndex(c echo.Context) error {
	data := TemplateData{
		IsAuthenticated: c.Get("isAuthenticated").(bool),
	}

	return c.Render(http.StatusOK, "index.html", data)
}

func (h *Handler) HandleSubmit(c echo.Context) error {
	h.logger.Info("Handling submit request")

	postalCode := c.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	// Server-side validation
	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.logger.Warn("Invalid postal code submitted", "postalCode", postalCode)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid postal code format")
	}

	mpFinder := services.NewMPFinder(h.client, h.logger)

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

	return templates.EmailTemplate.Execute(c.Response().Writer, data)
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
