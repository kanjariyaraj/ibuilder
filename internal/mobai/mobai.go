package mobai

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
)

type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateReconnecting
)

func (s ConnectionState) String() string {
	switch s {
	case StateDisconnected:
		return "Disconnected"
	case StateConnecting:
		return "Connecting"
	case StateConnected:
		return "Connected"
	case StateReconnecting:
		return "Reconnecting"
	default:
		return "Unknown"
	}
}

type DeviceInfo struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	OSVersion  string `json:"os_version"`
	UDID       string `json:"udid"`
	State      string `json:"state"`
	Battery    int    `json:"battery"`
	Storage    string `json:"storage"`
	Developer  bool   `json:"developer_mode"`
	Network    string `json:"network"`
}

type ConnectionStatus struct {
	State       ConnectionState `json:"state"`
	Device      *DeviceInfo     `json:"device,omitempty"`
	Latency     time.Duration   `json:"latency_ms"`
	ConnectedAt time.Time       `json:"connected_at,omitempty"`
	Error       string          `json:"error,omitempty"`
}

type Client struct {
	mu       sync.RWMutex
	cfg      *config.MobaiSettings
	log      *logger.Logger
	conn     net.Conn
	status   ConnectionStatus
	stopChan chan struct{}
}

func NewClient(cfg *config.MobaiSettings, log *logger.Logger) *Client {
	return &Client{
		cfg:      cfg,
		log:      log,
		status:   ConnectionStatus{State: StateDisconnected},
		stopChan: make(chan struct{}),
	}
}

func (c *Client) Config() *config.MobaiSettings {
	return c.cfg
}

func (c *Client) UpdateConfig(cfg *config.MobaiSettings) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cfg = cfg
}

func (c *Client) Status() ConnectionStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

func (c *Client) Connect() error {
	c.mu.Lock()
	if c.status.State == StateConnected {
		c.mu.Unlock()
		return fmt.Errorf("already connected")
	}
	c.status.State = StateConnecting
	c.mu.Unlock()

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
		timeout = 30 * time.Second
	}

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		c.mu.Lock()
		c.status = ConnectionStatus{
			State: StateDisconnected,
			Error: err.Error(),
		}
		c.mu.Unlock()
		c.log.Error("mobai connection failed", "addr", addr, "error", err)
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	device := &DeviceInfo{
		Name:      c.cfg.Device,
		Model:     "iPhone",
		OSVersion: "17.0",
		State:     "available",
		Battery:   85,
		Developer: true,
	}

	c.mu.Lock()
	c.conn = conn
	c.status = ConnectionStatus{
		State:       StateConnected,
		Device:      device,
		Latency:     5 * time.Millisecond,
		ConnectedAt: time.Now(),
	}
	c.mu.Unlock()

	c.log.Info("mobai connected", "addr", addr, "device", device.Name)

	if c.cfg.AutoReconnect {
		go c.reconnectLoop()
	}

	return nil
}

func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.stopChan:
	default:
		close(c.stopChan)
	}

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	c.status = ConnectionStatus{State: StateDisconnected}
	c.log.Info("mobai disconnected")
	return nil
}

func (c *Client) Ping() (time.Duration, error) {
	c.mu.RLock()
	if c.status.State != StateConnected {
		c.mu.RUnlock()
		return 0, fmt.Errorf("not connected")
	}
	c.mu.RUnlock()

	start := time.Now()
	_, err := fmt.Fprintf(c.conn, "ping\n")
	if err != nil {
		return 0, fmt.Errorf("ping failed: %w", err)
	}

	var resp string
	_, err = fmt.Fscan(c.conn, &resp)
	if err != nil {
		return 0, fmt.Errorf("ping response failed: %w", err)
	}

	latency := time.Since(start)

	c.mu.Lock()
	c.status.Latency = latency
	c.mu.Unlock()

	return latency, nil
}

func (c *Client) Doctor() HealthReport {
	report := HealthReport{
		Checks: []HealthCheck{},
	}

	report.Checks = append(report.Checks, c.checkConfig())
	report.Checks = append(report.Checks, c.checkConnectivity())
	report.Checks = append(report.Checks, c.checkDevice())

	for _, ch := range report.Checks {
		if ch.Status == StatusFailure {
			report.Healthy = false
			break
		}
	}
	if report.Healthy {
		for _, ch := range report.Checks {
			if ch.Status == StatusWarning {
				report.Healthy = false
				break
			}
		}
	}
	if len(report.Checks) == 0 {
		report.Healthy = true
	}

	return report
}

func (c *Client) reconnectLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			c.mu.RLock()
			if c.status.State == StateConnected {
				c.mu.RUnlock()
				continue
			}
			c.mu.RUnlock()

			c.mu.Lock()
			c.status.State = StateReconnecting
			c.mu.Unlock()

			c.log.Info("mobai reconnecting...")
			if err := c.Connect(); err != nil {
				c.log.Error("mobai reconnect failed", "error", err)
			}
		}
	}
}
