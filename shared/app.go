package shared

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
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
			logger, err := loggo.NewLogger(cfg.Log.File, cfg.GetLogLevel())
			if err != nil {
				return nil, fmt.Errorf("failed to create logger: %w", err)
			}
			return logger, nil
		},
		NewCustomFxLogger,
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
		fx.Annotate(
			NewFlashHandler,
			fx.As(new(FlashHandlerInterface)),
		),
	),
	ErrorModule,
)

// Provide a new database connection
func newDB(logger loggo.LoggerInterface, cfg *config.Config) (core.Interface, error) {
	logger.Info("Initializing database connection")

	retryConfig := dbconfig.NewDefaultRetryConfig()
	logger.Info("Attempting database connection with retry config",
		"maxAttempts", retryConfig.MaxAttempts,
		"initialInterval", retryConfig.InitialInterval)

	db, err := dbconfig.ConnectWithRetry(cfg, retryConfig, logger, &dbconfig.DefaultConnector{})
	if err != nil {
		logger.Error("Database connection failed", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Database connection successful")

	// List all files in migrations directory
	files, err := os.ReadDir(cfg.Server.MigrationsPath)
	if err != nil {
		logger.Error("Failed to read migrations directory", err)
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Log all SQL migration files found
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			logger.Info("Found SQL migration file", "filename", file.Name())

			// Optionally read and log migration content for debugging
			content, err := os.ReadFile(filepath.Join(cfg.Server.MigrationsPath, file.Name()))
			if err != nil {
				logger.Error("Failed to read migration file", err, "filename", file.Name())
			} else {
				logger.Info("Migration file content",
					"filename", file.Name(),
					"content", string(content))
			}
		}
	}

	decorated := &decorators.LoggingDecorator{
		Database: &core.DB{GormDB: db},
		Logger:   logger,
	}

	// Check database connection before migrations
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get underlying *sql.DB", err)
		return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Database ping failed", err)
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	logger.Info("Starting database migrations",
		"migrationsPath", cfg.Server.MigrationsPath,
		"databaseName", db.Migrator().CurrentDatabase())

	// Run migrations with detailed error capture
	if err := decorated.AutoMigrate(); err != nil {
		// Log the detailed error
		logger.Error("Database migration failed", err,
			"migrationsPath", cfg.Server.MigrationsPath,
			"databaseName", db.Migrator().CurrentDatabase(),
			"error", err.Error())

		// If it's a specific error type, log more details
		if migErr, ok := err.(interface{ Details() string }); ok {
			logger.Error("Migration error details", nil,
				"details", migErr.Details())
		}

		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return decorated, nil
}

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
		MailgunAPIKey: cfg.Email.MailgunKey,
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
