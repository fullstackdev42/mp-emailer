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
	"github.com/jonesrussell/mp-emailer/database"
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
	fx.Supply(fx.Hook{
		OnStart: func(context.Context) error {
			fmt.Println("🚀 Starting application...")
			return nil
		},
	}),
	fx.Provide(
		config.Load,
		provideVersionInfo,
		fx.Annotate(
			func(cfg *config.Config) (loggo.LoggerInterface, error) {
				fmt.Println("🔄 Initializing logger...")
				logger, err := loggo.NewLogger(cfg.Log.File, cfg.GetLogLevel())
				if err != nil {
					fmt.Printf("❌ Logger initialization failed: %v\n", err)
					return nil, fmt.Errorf("failed to create logger: %w", err)
				}
				fmt.Println("✅ Logger initialized successfully")
				return logger, nil
			},
		),
		NewCustomFxLogger,
		fx.Annotate(
			func() *echo.Echo {
				fmt.Println("🔄 Initializing Echo server...")
				e := echo.New()
				fmt.Println("✅ Echo server initialized")
				return e
			},
		),
		newSessionStore,
		validator.New,
		fx.Annotate(
			provideTemplates,
			fx.As(new(TemplateRendererInterface)),
		),
		fx.Annotate(
			provideEmailService,
			fx.As(new(email.Service)),
		),
		NewBaseHandler,
		NewGenericLoggingDecorator[LoggableService],
		fx.Annotate(
			NewFlashHandler,
			fx.As(new(FlashHandlerInterface)),
		),
		provideSessionCleaner,
		provideDatabaseService,
	),
	ErrorModule,
	fx.Invoke(
		func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					fmt.Println("🚀 Application initialization complete")
					return nil
				},
				OnStop: func(context.Context) error {
					fmt.Println("👋 Shutting down application...")
					return nil
				},
			})
		},
		startSessionCleaner,
	),
)

// Provide a new session store
func newSessionStore(cfg *config.Config, logger loggo.LoggerInterface) sessions.Store {
	logger.Debug("Initializing session store",
		"secure", cfg.App.Env == config.EnvProduction,
		"httpOnly", true,
		"maxAge", "7 days",
	)

	store := sessions.NewCookieStore([]byte(cfg.Auth.SessionSecret))

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   cfg.App.Env == config.EnvProduction,
		SameSite: http.SameSiteLaxMode,
	}

	logger.Debug("Session store initialized successfully")
	return store
}

// provideTemplates creates and configures the template renderer
func provideTemplates(store sessions.Store, cfg *config.Config, logger loggo.LoggerInterface) (TemplateRendererInterface, error) {
	logger.Debug("Initializing template renderer")

	tmpl := template.New("").Funcs(template.FuncMap{
		"hasPrefix": strings.HasPrefix,
		"safeHTML":  func(s string) template.HTML { return template.HTML(s) },
		"safeURL":   func(s string) template.URL { return template.URL(s) },
	})

	pattern := filepath.Join("web", "templates", "**", "*.gohtml")
	logger.Debug("Looking for templates", "pattern", pattern)

	templates, err := filepath.Glob(pattern)
	if err != nil {
		logger.Error("Failed to glob templates", err)
		return nil, fmt.Errorf("failed to glob templates: %w", err)
	}

	if len(templates) == 0 {
		logger.Error("No templates found", fmt.Errorf("no templates in %s", pattern))
		return nil, fmt.Errorf("no templates found in %s", pattern)
	}

	logger.Debug("Found templates", "count", len(templates))

	tmpl, err = tmpl.ParseFiles(templates...)
	if err != nil {
		logger.Error("Failed to parse templates", err)
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	logger.Debug("Templates parsed successfully")
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
func startSessionCleaner(lc fx.Lifecycle, cleaner *session.Cleaner, e *echo.Echo, logger loggo.LoggerInterface) {
	logger.Debug("Setting up session cleaner")

	e.Use(cleaner.Middleware())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting session cleanup routine")
			cleaner.StartCleanup(ctx)
			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("Stopping session cleanup routine")
			return nil
		},
	})
}

func provideDatabaseService(cfg *config.Config, logger loggo.LoggerInterface) (database.Database, error) {
	logger.Debug("Initializing database connection")
	ctx := context.Background()

	dbConfig := database.ConnectionConfig{
		DSN:                  cfg.DSN(),
		MaxRetries:           3,
		InitialInterval:      time.Second,
		MaxInterval:          time.Second * 10,
		MaxElapsedTime:       time.Minute,
		MultiplicationFactor: 2,
	}

	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		logger.Error("Failed to create database connection", err)
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}
	logger.Debug("Database connection established successfully")

	logger.Debug("Running database migrations")
	migConfig := database.MigrationConfig{
		DSN:            cfg.DSN(),
		MigrationsPath: "database/migrations",
		AllowDirty:     false,
	}

	if err := database.RunMigrations(migConfig); err != nil {
		logger.Error("Failed to run database migrations", err)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Debug("Database migrations completed successfully")

	return db, nil
}
