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
	logger loggo.LoggerInterface
	client api.ClientInterface
	store  *sessions.CookieStore
}

func NewHandler(logger loggo.LoggerInterface, client api.ClientInterface) *Handler {
	store := sessions.NewCookieStore([]byte("your-secret-key")) // Replace with a secure, randomly generated key
	return &Handler{logger: logger, client: client, store: store}
}

func (h *Handler) HandleIndex(c echo.Context) error {
	session, _ := h.store.Get(c.Request(), "session")
	user := session.Values["user"]

	data := struct {
		User interface{}
	}{
		User: user,
	}

	return templates.IndexTemplate.Execute(c.Response().Writer, data)
}

func (h *Handler) HandleSubmit(c echo.Context) error {
	h.logger.Info("Handling submit request")

	postalCode := c.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	// Server-side validation
	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.logger.Warn("Invalid postal code submitted", "postalCode", postalCode)
		return c.String(http.StatusBadRequest, "Invalid postal code format")
	}

	mpFinder := services.NewMPFinder(h.client, h.logger)

	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		h.logger.Error("Error finding MP", err)
		return c.String(http.StatusInternalServerError, "Error finding MP")
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
		h.logger.Error("Error binding request", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	return c.JSON(http.StatusOK, req)
}

func composeEmail(mp models.Representative) string {
	return fmt.Sprintf("Dear %s,\n\nThis is a sample email content.\n\nBest regards,\nYour constituent", mp.Name)
}

func (h *Handler) HandleLogin(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return templates.LoginTemplate.Execute(c.Response().Writer, nil)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	// TODO: Implement proper user authentication
	if username == "admin" && password == "password" {
		session, _ := h.store.Get(c.Request(), "session")
		session.Values["user"] = username
		session.Save(c.Request(), c.Response().Writer)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.String(http.StatusUnauthorized, "Invalid credentials")
}

func (h *Handler) HandleLogout(c echo.Context) error {
	session, _ := h.store.Get(c.Request(), "session")
	session.Values["user"] = nil
	session.Save(c.Request(), c.Response().Writer)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := h.store.Get(c.Request(), "session")
		user := session.Values["user"]
		if user == nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		return next(c)
	}
}
