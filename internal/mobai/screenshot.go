package mobai

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

type ScreenshotResult struct {
	Path      string    `json:"path"`
	Size      int64     `json:"size_bytes"`
	CapturedAt time.Time `json:"captured_at"`
	Device    string    `json:"device"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
}

func (c *Client) CaptureScreenshot(outputDir string) (*ScreenshotResult, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to MobAI")
	}

	if outputDir == "" {
		outputDir = "screenshots"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create screenshot directory: %w", err)
	}

	deviceName := "unknown"
	c.mu.RLock()
	if c.status.Device != nil {
		deviceName = c.status.Device.Name
	}
	c.mu.RUnlock()

	timestamp := c.Timestamp()
	filename := fmt.Sprintf("screenshot_%s_%s.png", sanitizeFilename(deviceName), timestamp)
	path := filepath.Join(outputDir, filename)

	c.log.Info("capturing screenshot", "device", deviceName)

	img := generateMockScreenshot()
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create screenshot file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return nil, fmt.Errorf("failed to encode screenshot: %w", err)
	}

	info, _ := os.Stat(path)

	result := &ScreenshotResult{
		Path:       path,
		Size:       info.Size(),
		CapturedAt: time.Now(),
		Device:     deviceName,
		Width:      1179,
		Height:     2556,
	}

	c.log.Info("screenshot captured", "path", path, "size", result.Size)
	return result, nil
}

func sanitizeFilename(name string) string {
	clean := make([]byte, 0, len(name))
	for _, b := range []byte(name) {
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '-' || b == '_' {
			clean = append(clean, b)
		} else if b == ' ' {
			clean = append(clean, '_')
		}
	}
	if len(clean) == 0 {
		return "device"
	}
	return string(clean)
}

func generateMockScreenshot() image.Image {
	width, height := 1179, 2556
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	bgColor := color.RGBA{30, 30, 40, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	statusBarColor := color.RGBA{20, 20, 30, 255}
	for x := 0; x < width; x++ {
		for y := 0; y < 60; y++ {
			img.Set(x, y, statusBarColor)
		}
	}

	accentColor := color.RGBA{0, 122, 255, 255}
	centerX, centerY := width/2, height/2
	for dx := -50; dx <= 50; dx++ {
		for dy := -50; dy <= 50; dy++ {
			px, py := centerX+dx, centerY+dy
			if px >= 0 && px < width && py >= 0 && py < height {
				img.Set(px, py, accentColor)
			}
		}
	}

	return img
}
