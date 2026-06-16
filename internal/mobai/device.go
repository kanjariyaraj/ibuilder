package mobai

import (
	"fmt"
	"time"
)

var mockDevices = []DeviceInfo{
	{
		Name:      "Raj's iPhone",
		Model:     "iPhone 15 Pro",
		OSVersion: "17.4",
		UDID:      "00008110-xxxxxxxxxxxx",
		State:     "available",
		Battery:   85,
		Storage:   "128GB / 256GB",
		Developer: true,
		Network:   "Wi-Fi",
	},
	{
		Name:      "Test iPhone 14",
		Model:     "iPhone 14",
		OSVersion: "17.2",
		UDID:      "00008120-yyyyyyyyyyyy",
		State:     "available",
		Battery:   72,
		Storage:   "64GB / 128GB",
		Developer: true,
		Network:   "USB",
	},
}

func (c *Client) ListDevices() ([]DeviceInfo, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to MobAI")
	}

	return mockDevices, nil
}

func (c *Client) DeviceInfo(udid string) (*DeviceInfo, error) {
	devices, err := c.ListDevices()
	if err != nil {
		return nil, err
	}

	if udid == "" && len(devices) > 0 {
		return &devices[0], nil
	}

	for _, d := range devices {
		if d.UDID == udid {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("device %s not found", udid)
}

func FormatDeviceInfo(d *DeviceInfo) string {
	return fmt.Sprintf(`Device Information:
  Name:           %s
  Model:          %s
  iOS Version:    %s
  UDID:           %s
  State:          %s
  Battery:        %d%%
  Storage:        %s
  Developer Mode: %v
  Network:        %s
`, d.Name, d.Model, d.OSVersion, d.UDID, d.State, d.Battery, d.Storage, d.Developer, d.Network)
}

func PrintDeviceTable(devices []DeviceInfo) {
	fmt.Printf("%-25s %-18s %-10s %-12s %-8s %-10s\n",
		"Name", "Model", "iOS", "State", "Battery", "UDID")
	fmt.Println("----------------------------------------------------------------------")
	for _, d := range devices {
		fmt.Printf("%-25s %-18s %-10s %-12s %-8d %-10s\n",
			truncate(d.Name, 24),
			d.Model,
			d.OSVersion,
			d.State,
			d.Battery,
			truncate(d.UDID, 10),
		)
	}
}

func (c *Client) Timestamp() string {
	return time.Now().Format("20060102_150405")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
