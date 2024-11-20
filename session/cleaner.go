package session

import (
	"context"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type Cleaner struct {
	store           sessions.Store
	cleanupInterval time.Duration
	maxAge          int
	logger          loggo.LoggerInterface
}

func NewCleaner(store sessions.Store, cleanupInterval time.Duration, maxAge int, logger loggo.LoggerInterface) *Cleaner {
	return &Cleaner{
		store:           store,
		cleanupInterval: cleanupInterval,
		maxAge:          maxAge,
		logger:          logger,
	}
}

// StartCleanup initiates the periodic session cleanup
func (sc *Cleaner) StartCleanup(ctx context.Context) {
	ticker := time.NewTicker(sc.cleanupInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				sc.cleanup()
			}
		}
	}()
}

// cleanup performs the actual session cleanup
func (sc *Cleaner) cleanup() {
	// For CookieStore, we don't need to check the store type
	// as sessions are managed by the browser via cookie expiration
	now := time.Now()
	threshold := now.Add(-time.Duration(sc.maxAge) * time.Second)

	sc.logger.Info("Session cleanup check",
		"maxAge", sc.maxAge,
		"threshold", threshold,
		"cleanupInterval", sc.cleanupInterval)
}

// Middleware returns an Echo middleware function
func (sc *Cleaner) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the session
			sess, err := sc.store.Get(c.Request(), "session")
			if err != nil {
				sc.logger.Error("Error getting session", err)
				return next(c)
			}

			// Check if session is expired
			if isExpired(sess) {
				// Delete the session
				sess.Options.MaxAge = -1
				err = sess.Save(c.Request(), c.Response().Writer)
				if err != nil {
					sc.logger.Error("Error deleting expired session", err)
				}
			}

			return next(c)
		}
	}
}

// isExpired checks if a session is expired
func isExpired(s *sessions.Session) bool {
	lastAccessed, ok := s.Values["last_accessed"].(time.Time)
	if !ok {
		return true
	}

	// Check if session has exceeded max age
	maxAge := s.Options.MaxAge
	if maxAge > 0 {
		return time.Since(lastAccessed) > time.Duration(maxAge)*time.Second
	}

	return false
}
