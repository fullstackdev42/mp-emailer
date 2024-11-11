package shared

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/fullstackdev42/mp-emailer/config"
	dbconfig "github.com/fullstackdev42/mp-emailer/database/config"
	"github.com/fullstackdev42/mp-emailer/database/core"
	"github.com/fullstackdev42/mp-emailer/database/decorators"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// App provides the shared application modules
//
//nolint:gochecknoglobals
var App = fx.Options(
	fx.Provide(
		config.Load,
		func(cfg *config.Config) (loggo.LoggerInterface, error) {
			logger, err := loggo.NewLogger(cfg.LogFile, cfg.GetLogLevel())
			if err != nil {
				return nil, fmt.Errorf("failed to create logger: %w", err)
			}
			return logger, nil
		},
		fx.Annotate(
			newDB,
			fx.As(new(core.Interface)),
		),
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
		NewFlashHandler,
	),
	ErrorModule,
)

// Provide a new database connection
func newDB(logger loggo.LoggerInterface, cfg *config.Config) (core.Interface, error) {
	logger.Info("Initializing database connection")

	// Use the proper retry configuration
	retryConfig := dbconfig.NewDefaultRetryConfig()
	db, err := dbconfig.ConnectWithRetry(cfg, retryConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Wrap the DB in a decorator that implements the correct interface
	decorated := &decorators.LoggingDecorator{
		Database: &core.DB{GormDB: db},
		Logger:   logger,
	}

	// Ensure database migrations are run
	if err := decorated.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return decorated, nil
}

// Provide a new session store
func newSessionStore(cfg *config.Config) sessions.Store {
	return sessions.NewCookieStore([]byte(cfg.SessionSecret))
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
		Provider:      email.Provider(cfg.EmailProvider),
		SMTPHost:      cfg.SMTPHost,
		SMTPPort:      cfg.SMTPPort,
		SMTPUsername:  cfg.SMTPUsername,
		SMTPPassword:  cfg.SMTPPassword,
		SMTPFrom:      cfg.SMTPFrom,
		MailgunDomain: cfg.MailgunDomain,
		MailgunAPIKey: cfg.MailgunAPIKey,
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
