package shared

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/email"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/jonesrussell/mp-emailer/version"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// App provides the shared application modules
//
//nolint:gochecknoglobals
var App = fx.Options(
	fx.Provide(
		config.Load,
		provideVersionInfo,
		func(cfg *config.Config) (loggo.LoggerInterface, error) {
			logger, err := loggo.NewLogger(cfg.Log.File, cfg.GetLogLevel())
			if err != nil {
				return nil, fmt.Errorf("failed to create logger: %w", err)
			}
			return logger, nil
		},
		NewCustomFxLogger,
		echo.New,
		newSessionStore,
		validator.New,
		fx.Annotate(
			provideTemplates,
			fx.As(new(TemplateRendererInterface)),
		),

		provideEmailService,
		NewBaseHandler,
		NewGenericLoggingDecorator[LoggableService],
		fx.Annotate(
			NewFlashHandler,
			fx.As(new(FlashHandlerInterface)),
		),
		provideSessionCleaner,
	),
	ErrorModule,
	fx.Invoke(
		startSessionCleaner,
	),
)

// Provide a new session store
func newSessionStore(cfg *config.Config) sessions.Store {
	store := sessions.NewCookieStore([]byte(cfg.Auth.SessionSecret))

	// Configure secure cookie options
	store.Options = &sessions.Options{
		Path:     "/",                     // Cookie available for entire site
		MaxAge:   86400 * 7,               // 7 days in seconds
		HttpOnly: true,                    // Prevent XSS by making cookie inaccessible to JS
		Secure:   true,                    // Only send cookie over HTTPS
		SameSite: http.SameSiteStrictMode, // Strict SameSite policy for CSRF protection
	}

	return store
}

// provideTemplates creates and configures the template renderer
func provideTemplates(store sessions.Store, cfg *config.Config) (TemplateRendererInterface, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"hasPrefix": strings.HasPrefix,
		"safeHTML":  func(s string) template.HTML { return template.HTML(s) },
		"safeURL":   func(s string) template.URL { return template.URL(s) },
	})

	pattern := filepath.Join("web", "templates", "**", "*.gohtml")
	templates, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob templates: %w", err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates found in %s", pattern)
	}

	tmpl, err = tmpl.ParseFiles(templates...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return NewCustomTemplateRenderer(tmpl, store, cfg), nil
}

// provideEmailService creates a new email service based on the configuration
func provideEmailService(cfg *config.Config, logger loggo.LoggerInterface) (email.Service, error) {
	emailConfig := email.Config{
		Provider:      email.Provider(cfg.Email.Provider),
		SMTPHost:      cfg.Email.SMTP.Host,
		SMTPPort:      cfg.Email.SMTP.Port,
		SMTPUsername:  cfg.Email.SMTP.Username,
		SMTPPassword:  cfg.Email.SMTP.Password,
		SMTPFrom:      cfg.Email.SMTP.From,
		MailgunDomain: cfg.Email.MailgunDomain,
		MailgunAPIKey: cfg.Email.MailgunAPIKey,
	}

	emailService, err := email.NewEmailService(email.Params{
		Config: emailConfig,
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create email service: %w", err)
	}

	return emailService, nil
}

func provideVersionInfo() version.Info {
	return version.Get()
}

// Add provider function for SessionCleaner
func provideSessionCleaner(store sessions.Store, cfg *config.Config, logger loggo.LoggerInterface) *session.Cleaner {
	return session.NewCleaner(
		store,
		15*time.Minute, // cleanup interval
		cfg.Auth.SessionMaxAge,
		logger,
	)
}

// Add startup function
func startSessionCleaner(lc fx.Lifecycle, cleaner *session.Cleaner, e *echo.Echo) {
	// Add the cleanup middleware to Echo
	e.Use(cleaner.Middleware())

	// Start the cleanup routine
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			cleaner.StartCleanup(ctx)
			return nil
		},
	})
}
