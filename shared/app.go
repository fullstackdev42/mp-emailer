package shared

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/database"
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
			fx.As(new(database.Interface)),
		),
		echo.New,
		newSessionStore,
		validator.New,
		fx.Annotate(
			ProvideTemplates,
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
func newDB(logger loggo.LoggerInterface, cfg *config.Config) (database.Interface, error) {
	logger.Info("Initializing database connection")
	dsn := cfg.DatabaseDSN()

	db, err := connectToDB(dsn, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after multiple attempts: %w", err)
	}

	// Wrap the base DB with the logging decorator
	decorated := database.NewLoggingDBDecorator(db, logger)
	return decorated, nil
}

// connectToDB attempts to connect to the database with retries
func connectToDB(dsn string, logger loggo.LoggerInterface) (database.Interface, error) {
	var err error
	for retries := 5; retries > 0; retries-- {
		db, err := database.NewDB(dsn, logger)
		if err == nil {
			return db, nil
		}
		logger.Warn("Failed to connect to database, retrying...", "error", err)
		time.Sleep(5 * time.Second)
	}
	return nil, err
}

// Provide a new session store
func newSessionStore(cfg *config.Config) sessions.Store {
	return sessions.NewCookieStore([]byte(cfg.SessionSecret))
}

// ProvideTemplates creates and configures the template renderer
func ProvideTemplates(store sessions.Store) (TemplateRendererInterface, error) {
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

	return NewCustomTemplateRenderer(tmpl, store), nil
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

	emailService, err := email.NewEmailService(emailConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create email service: %w", err)
	}

	return emailService, nil
}
