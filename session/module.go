package session

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/logger"
	"go.uber.org/fx"
)

// Module provides all session-related dependencies
//
//nolint:gochecknoglobals
var Module = fx.Module("session",
	fx.Provide(
		NewDefaultOptions,
		// Keep only this provider for StoreProvider
		fx.Annotate(
			func(store Store) StoreProvider {
				return NewStoreProvider(store)
			},
			fx.As(new(StoreProvider)),
		),
		// Keep the secure store provider
		fx.Annotate(
			func(options Options) (Store, error) {
				store := sessions.NewCookieStore(options.SecurityKey)
				return NewSecureStore(store, options)
			},
		),
	),
)

func NewDefaultOptions(cfg *config.Config, log logger.Interface) (Options, error) {
	// Get the session secret from config
	secret := cfg.Auth.SessionSecret

	// Ensure the secret is base64 encoded
	if _, err := base64.StdEncoding.DecodeString(secret); err != nil {
		return Options{}, fmt.Errorf("session secret must be valid base64: %w", err)
	}

	// Decode the base64 secret key
	decodedKey, err := base64.StdEncoding.DecodeString(cfg.Auth.SessionSecret)
	if err != nil {
		return Options{}, fmt.Errorf("session secret must be valid base64: %w", err)
	}

	// Check the decoded key length
	if len(decodedKey) != 32 {
		log.Error("Invalid session secret length", fmt.Errorf("decoded session secret must be exactly 32 bytes (current length: %d)", len(decodedKey)))
		return Options{}, fmt.Errorf("decoded session secret must be exactly 32 bytes (got %d)", len(decodedKey))
	}

	return Options{
		MaxAge:          cfg.Auth.SessionMaxAge,
		CleanupInterval: 15 * time.Minute,
		SecurityKey:     decodedKey,
		CookieName:      cfg.Auth.SessionName,
		Domain:          cfg.App.Domain,
		Secure:          cfg.App.Env == "production",
		HTTPOnly:        true,
		SameSite:        http.SameSiteLaxMode,
		Path:            "/",
		MaxLength:       4096, // 4KB
		KeyPrefix:       "sess_",
	}, nil
}
