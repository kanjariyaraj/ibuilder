package mobai

import (
	"fmt"
	"time"
)

type ReconnectStats struct {
	Attempts     int           `json:"attempts"`
	Successes    int           `json:"successes"`
	Failures     int           `json:"failures"`
	LastAttempt  time.Time     `json:"last_attempt"`
	LastSuccess  time.Time     `json:"last_success"`
	TotalTime    time.Duration `json:"total_downtime"`
}

func (c *Client) Reconnect() error {
	c.mu.Lock()
	if c.status.State == StateConnected {
		c.mu.Unlock()
		return nil
	}
	c.status.State = StateReconnecting
	c.mu.Unlock()

	c.log.Info("attempting manual reconnect")

	if err := c.Disconnect(); err != nil {
		c.log.Error("disconnect during reconnect failed", "error", err)
	}

	if err := c.Connect(); err != nil {
		c.log.Error("manual reconnect failed", "error", err)
		return fmt.Errorf("reconnect failed: %w", err)
	}

	c.log.Info("manual reconnect successful")
	return nil
}

func (c *Client) SessionRestore() error {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if connected {
		return nil
	}

	c.log.Info("attempting session restoration")

	cfg := c.cfg
	if cfg.Host == "" || cfg.Port == 0 {
		return fmt.Errorf("cannot restore session: incomplete configuration")
	}

	if err := c.Connect(); err != nil {
		return fmt.Errorf("session restoration failed: %w", err)
	}

	return nil
}

func (c *Client) IsAutoReconnectEnabled() bool {
	return c.cfg.AutoReconnect
}

func (c *Client) SetAutoReconnect(enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cfg.AutoReconnect = enabled
}
