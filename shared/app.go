package shared

import (
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/database"
	"github.com/jonesrussell/mp-emailer/email"
	"github.com/jonesrussell/mp-emailer/logger"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/jonesrussell/mp-emailer/version"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// App is the main application module
//
//nolint:gochecknoglobals
var App = fx.Options(
	fx.Provide(
		config.Load,
		provideVersionInfo,
		fx.Annotate(
			provideLogger,
			fx.As(new(logger.Interface)),
		),
		echo.New,
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
		provideDatabaseService,
		fx.Annotate(
			func(log logger.Interface, options session.Options) (session.Manager, error) {
				store := sessions.NewCookieStore(options.SecurityKey)
				secureStore, err := session.NewSecureStore(store, options)
				if err != nil {
					return nil, err
				}

				manager, err := session.NewManager(secureStore, log, options)
				if err != nil {
					return nil, fmt.Errorf("failed to create session manager: %w", err)
				}

				return manager, nil
			},
			fx.As(new(session.Manager)),
		),
	),
	ErrorModule,
)

func provideLogger(cfg *config.Config) (logger.Interface, error) {
	logConfig := &logger.Config{
		Level:       cfg.Log.Level,
		Format:      cfg.Log.Format,
		OutputPath:  cfg.Log.File,
		Development: cfg.App.Env == config.EnvDevelopment,
	}

	if err := logger.Initialize(logConfig); err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return logger.GetLogger(), nil
}

func provideVersionInfo() version.Info {
	return version.Get()
}

func provideDatabaseService(cfg *config.Config) (database.Database, error) {
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
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	migConfig := database.MigrationConfig{
		DSN:            cfg.DSN(),
		MigrationsPath: "database/migrations",
		AllowDirty:     false,
	}

	if err := database.RunMigrations(migConfig); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func provideTemplates(manager session.Manager, cfg *config.Config) (TemplateRendererInterface, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"hasPrefix": strings.HasPrefix,
		"safeHTML":  func(s string) template.HTML { return template.HTML(s) },
		"safeURL":   func(s string) template.URL { return template.URL(s) },
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
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

	return NewCustomTemplateRenderer(tmpl, manager, cfg), nil
}

func provideEmailService(cfg *config.Config, log logger.Interface) (email.Service, error) {
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
		Logger: log,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create email service: %w", err)
	}

	return emailService, nil
}
