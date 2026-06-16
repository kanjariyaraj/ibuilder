package reactnative

import (
	"fmt"
	"os/exec"
	"strings"
)

type InstallResult struct {
	Success  bool   `json:"success"`
	Action   string `json:"action"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
}

func (s *Session) InstallLatest() (*InstallResult, error) {
	s.log.Info("installing latest react native build")

	deviceID := s.deviceID
	args := []string{"react-native", "run-ios", "--mode", "Release"}
	if deviceID != "" {
		args = append(args, "--device", deviceID)
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		result := &InstallResult{
			Success: false,
			Action:  "install_latest",
			Output:  outputStr,
			Error:   err.Error(),
		}
		return result, fmt.Errorf("install failed: %s", outputStr)
	}

	result := &InstallResult{
		Success: true,
		Action:  "install_latest",
		Output:  outputStr,
	}

	s.log.Info("latest react native app installed")
	return result, nil
}

func (s *Session) InstallSpecificBuild(artifactPath string) (*InstallResult, error) {
	s.log.Info("installing specific react native build", "artifact", artifactPath)

	if artifactPath == "" {
		return nil, fmt.Errorf("artifact path is required")
	}

	args := []string{"react-native", "run-ios", "--binary", artifactPath}
	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		result := &InstallResult{
			Success: false,
			Action:  "install_specific",
			Output:  outputStr,
			Error:   err.Error(),
		}
		return result, fmt.Errorf("install failed: %s", outputStr)
	}

	result := &InstallResult{
		Success: true,
		Action:  "install_specific",
		Output:  outputStr,
	}

	s.log.Info("specific build installed", "artifact", artifactPath)
	return result, nil
}

func (s *Session) VerifyInstallation() error {
	s.log.Info("verifying installation")

	args := []string{"react-native", "list-devices"}
	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("installation verification failed: %s", strings.TrimSpace(string(output)))
	}

	if !strings.Contains(string(output), "device") {
		return fmt.Errorf("no device detected for verification")
	}

	s.log.Info("installation verified")
	return nil
}
