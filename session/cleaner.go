package session

import (
	"context"
	"time"

	"github.com/jonesrussell/mp-emailer/logger"
)

type Cleaner struct {
	store    Store
	interval time.Duration
	maxAge   int
	logger   logger.Interface
	cancel   context.CancelFunc
}

func NewCleaner(store Store, interval time.Duration, maxAge int, log logger.Interface) *Cleaner {
	return &Cleaner{
		store:    store,
		interval: interval,
		maxAge:   maxAge,
		logger:   log,
	}
}

func (c *Cleaner) StartCleanup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.logger.Debug("Starting session cleanup routine", "interval", c.interval, "maxAge", c.maxAge)

	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.logger.Debug("Running session cleanup")
				threshold := time.Now().Add(-time.Duration(c.maxAge) * time.Second)
				if err := c.store.Cleanup(threshold); err != nil {
					c.logger.Error("Session cleanup failed", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *Cleaner) StopCleanup() {
	if c.cancel != nil {
		c.cancel()
	}
}
