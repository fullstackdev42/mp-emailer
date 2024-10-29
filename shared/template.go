package shared

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// CustomTemplateRenderer is a custom renderer for Echo
type CustomTemplateRenderer struct {
	templates *template.Template
}

// NewCustomTemplateRenderer creates a new CustomTemplateRenderer
func NewCustomTemplateRenderer(templates *template.Template) *CustomTemplateRenderer {
	return &CustomTemplateRenderer{templates: templates}
}

// Render renders a template document
func (t *CustomTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Get the session store from the context
	store, ok := c.Get("store").(sessions.Store)
	if !ok {
		return fmt.Errorf("could not get session store from context")
	}

	// Use the same session name as in the handlers
	session, err := store.Get(c.Request(), "mpe")
	if err != nil {
		c.Logger().Errorf("failed to get session: %v", err)
		return err
	}

	// Get and clear flash messages
	flashes := session.Flashes("messages")
	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		c.Logger().Errorf("failed to save session: %v", err)
		return err
	}

	viewData := map[string]interface{}{
		"Data":      data,
		"RequestID": c.Response().Header().Get(echo.HeaderXRequestID),
		"PageName":  name,
		"Messages":  flashes,
	}

	var content bytes.Buffer
	if err := t.templates.ExecuteTemplate(&content, name, viewData); err != nil {
		c.Logger().Errorf("failed to execute template %s: %v", name, err)
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}
	viewData["TemplateContent"] = template.HTML(content.String())
	return t.templates.ExecuteTemplate(w, "app", viewData)
}

func (t *CustomTemplateRenderer) RenderPage(c echo.Context, template string, pageData PageData, logger loggo.LoggerInterface, errorHandler *ErrorHandler) error {
	err := t.Render(c.Response(), template, pageData, c)
	if err != nil {
		logger.Error("Failed to render template", err)
		return errorHandler.HandleHTTPError(c, err, "Failed to render page", http.StatusInternalServerError)
	}
	return nil
}
