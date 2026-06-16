package mobai

import (
	"fmt"
	"net"
	"time"
)

type HealthStatus string

const (
	StatusHealthy  HealthStatus = "HEALTHY"
	StatusWarning  HealthStatus = "WARNING"
	StatusFailure  HealthStatus = "FAILURE"
)

type HealthCheck struct {
	Name    string       `json:"name"`
	Status  HealthStatus `json:"status"`
	Message string       `json:"message"`
	Suggest string       `json:"suggest,omitempty"`
}

type HealthReport struct {
	Healthy bool          `json:"healthy"`
	Checks  []HealthCheck `json:"checks"`
}

func (c *Client) checkConfig() HealthCheck {
	cfg := c.cfg

	if cfg.Host == "" {
		return HealthCheck{
			Name:    "Configuration",
			Status:  StatusWarning,
			Message: "Host not configured, using default (localhost)",
			Suggest: "Set mobai.host in builder.json",
		}
	}

	port := cfg.Port
	if port == 0 {
		port = 12345
	}

	return HealthCheck{
		Name:    "Configuration",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("Configured: %s:%d", cfg.Host, port),
	}
}

func (c *Client) checkConnectivity() HealthCheck {
	host := c.cfg.Host
	port := c.cfg.Port
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 12345
	}

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	timeout := time.Duration(c.cfg.ConnectionTimeout) * time.Second
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return HealthCheck{
			Name:    "Connectivity",
			Status:  StatusFailure,
			Message: fmt.Sprintf("Cannot reach %s: %v", addr, err),
			Suggest: "Ensure MobAI service is running on the host",
		}
	}
	conn.Close()

	return HealthCheck{
		Name:    "Connectivity",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("Reachable: %s (%dms)", addr, timeout.Milliseconds()),
	}
}

func (c *Client) checkDevice() HealthCheck {
	c.mu.RLock()
	state := c.status.State
	device := c.status.Device
	c.mu.RUnlock()

	if state != StateConnected {
		return HealthCheck{
			Name:    "Device",
			Status:  StatusWarning,
			Message: "No device connected",
			Suggest: "Run 'builder mobai connect' first",
		}
	}

	if device == nil {
		return HealthCheck{
			Name:    "Device",
			Status:  StatusFailure,
			Message: "Device info not available",
			Suggest: "Reconnect and try again",
		}
	}

	return HealthCheck{
		Name:    "Device",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("%s (%s) - iOS %s", device.Name, device.Model, device.OSVersion),
	}
}
