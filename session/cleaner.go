package session

import (
	"context"
	"sync"
	"time"

	"github.com/jonesrussell/loggo"
)

type Cleaner struct {
	store         Store
	interval      time.Duration
	maxAge        int
	logger        loggo.LoggerInterface
	cleanupTicker *time.Ticker
	cleanupDone   chan struct{}
	cleanupWG     sync.WaitGroup
}

func NewCleaner(store Store, interval time.Duration, maxAge int, logger loggo.LoggerInterface) *Cleaner {
	return &Cleaner{
		store:       store,
		interval:    interval,
		maxAge:      maxAge,
		logger:      logger,
		cleanupDone: make(chan struct{}),
	}
}

func (c *Cleaner) StartCleanup(ctx context.Context) {
	c.logger.Debug("Starting session cleanup routine",
		"interval", c.interval,
		"maxAge", c.maxAge)

	c.cleanupTicker = time.NewTicker(c.interval)
	c.cleanupWG.Add(1)

	go func() {
		defer c.cleanupWG.Done()
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("Context cancelled, stopping cleanup routine")
				return
			case <-c.cleanupDone:
				c.logger.Debug("Cleanup routine stopped")
				return
			case <-c.cleanupTicker.C:
				c.cleanup()
			}
		}
	}()
}

func (c *Cleaner) StopCleanup() error {
	c.logger.Debug("Stopping session cleanup routine")

	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
		close(c.cleanupDone)
		c.cleanupWG.Wait()
	}
	return nil
}

func (c *Cleaner) cleanup() {
	c.logger.Debug("Running session cleanup")

	threshold := time.Now().Add(-time.Duration(c.maxAge) * time.Second)

	if store, ok := c.store.(interface {
		Cleanup(threshold time.Time) error
	}); ok {
		if err := store.Cleanup(threshold); err != nil {
			c.logger.Error("Session cleanup failed", err)
		}
	}
}
