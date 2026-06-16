package reactnative

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type InstallResult struct {
	Success   bool      `json:"success"`
	Artifact  string    `json:"artifact,omitempty"`
	Device    string    `json:"device,omitempty"`
	Installed time.Time `json:"installed_at"`
	Error     string    `json:"error,omitempty"`
}

func (s *Session) InstallLatest() (*InstallResult, error) {
	s.log.Info("installing latest RN build")

	args := []string{"react-native", "run-ios", "--no-launch", "--json"}
	if s.deviceID != "" {
		args = append(args, "--device", s.deviceID)
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("install failed: %s", strings.TrimSpace(string(output)))
	}

	result := &InstallResult{
		Success:   true,
		Artifact:  "latest",
		Device:    s.deviceID,
		Installed: time.Now(),
	}

	s.log.Info("RN app installed", "device", s.deviceID)
	return result, nil
}

func (s *Session) InstallArtifact(artifactPath string) (*InstallResult, error) {
	s.log.Info("installing RN build artifact", "path", artifactPath)

	args := []string{"react-native", "run-ios", "--no-launch", "--binary", artifactPath}
	if s.deviceID != "" {
		args = append(args, "--device", s.deviceID)
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("artifact install failed: %s", strings.TrimSpace(string(output)))
	}

	result := &InstallResult{
		Success:   true,
		Artifact:  artifactPath,
		Device:    s.deviceID,
		Installed: time.Now(),
	}

	s.log.Info("RN artifact installed", "device", s.deviceID, "artifact", artifactPath)
	return result, nil
}

func (s *Session) VerifyInstall() (bool, error) {
	s.mu.RLock()
	deviceID := s.deviceID
	s.mu.RUnlock()

	args := []string{"react-native", "list-devices", "--json"}
	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to verify install: %w", err)
	}

	installed := false
	for _, line := range strings.Split(string(output), "\n") {
		if deviceID != "" && strings.Contains(line, deviceID) {
			installed = true
			break
		}
		if strings.Contains(line, "device") {
			installed = true
		}
	}

	return installed, nil
}
