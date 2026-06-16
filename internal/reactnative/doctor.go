package reactnative

import (
	"fmt"
	"os/exec"
	"strings"
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

type DoctorReport struct {
	Healthy bool          `json:"healthy"`
	Checks  []HealthCheck `json:"checks"`
}

func (s *Session) Doctor() DoctorReport {
	report := DoctorReport{
		Checks: []HealthCheck{},
	}

	report.Checks = append(report.Checks, s.checkNode())
	report.Checks = append(report.Checks, s.checkNPM())
	report.Checks = append(report.Checks, s.checkProject())
	report.Checks = append(report.Checks, s.checkDependencies())
	report.Checks = append(report.Checks, s.checkMetro())
	report.Checks = append(report.Checks, s.checkDevices())

	for _, ch := range report.Checks {
		if ch.Status == StatusFailure {
			report.Healthy = false
			return report
		}
	}
	report.Healthy = true
	return report
}

func (s *Session) checkNode() HealthCheck {
	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return HealthCheck{
			Name:    "Node.js",
			Status:  StatusFailure,
			Message: "Node.js is not installed or not in PATH",
			Suggest: "Install Node.js from https://nodejs.org",
		}
	}

	return HealthCheck{
		Name:    "Node.js",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("Node %s installed", strings.TrimSpace(string(output))),
	}
}

func (s *Session) checkNPM() HealthCheck {
	cmd := exec.Command("npm", "--version")
	output, err := cmd.Output()
	if err != nil {
		return HealthCheck{
			Name:    "npm",
			Status:  StatusFailure,
			Message: "npm is not installed or not in PATH",
			Suggest: "npm is bundled with Node.js — reinstall Node.js",
		}
	}

	return HealthCheck{
		Name:    "npm",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("npm %s installed", strings.TrimSpace(string(output))),
	}
}

func (s *Session) checkProject() HealthCheck {
	if s.projectDir == "" {
		return HealthCheck{
			Name:    "React Native Project",
			Status:  StatusWarning,
			Message: "No project directory set",
			Suggest: "Run 'builder rn dev' from within an RN project directory",
		}
	}

	valid, err := s.DetectRNProject(s.projectDir)
	if err != nil {
		return HealthCheck{
			Name:    "React Native Project",
			Status:  StatusFailure,
			Message: fmt.Sprintf("Invalid project: %v", err),
			Suggest: "Ensure the directory contains a React Native project with package.json containing react-native and ios/",
		}
	}
	if !valid {
		return HealthCheck{
			Name:    "React Native Project",
			Status:  StatusFailure,
			Message: "Not a valid React Native project",
			Suggest: "Run 'npx react-native init' or navigate to an RN project directory",
		}
	}

	return HealthCheck{
		Name:    "React Native Project",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("Valid React Native project in %s", s.projectDir),
	}
}

func (s *Session) checkDependencies() HealthCheck {
	if s.projectDir == "" {
		return HealthCheck{
			Name:    "Dependencies",
			Status:  StatusWarning,
			Message: "No project directory set",
		}
	}

	nodeModules := s.projectDir + "/node_modules"
	if _, err := exec.Command("test", "-d", nodeModules).CombinedOutput(); err != nil {
		return HealthCheck{
			Name:    "Dependencies",
			Status:  StatusWarning,
			Message: "node_modules not found — dependencies not installed",
			Suggest: "Run 'npm install' in your project directory",
		}
	}

	return HealthCheck{
		Name:    "Dependencies",
		Status:  StatusHealthy,
		Message: "Dependencies installed",
	}
}

func (s *Session) checkMetro() HealthCheck {
	cmd := exec.Command("npx", "--yes", "react-native", "--help")
	if err := cmd.Run(); err != nil {
		return HealthCheck{
			Name:    "Metro Bundler",
			Status:  StatusWarning,
			Message: "npx react-native not available",
			Suggest: "Ensure react-native is installed in your project dependencies",
		}
	}

	return HealthCheck{
		Name:    "Metro Bundler",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("Metro available (port %d)", s.metroPort),
	}
}

func (s *Session) checkDevices() HealthCheck {
	args := []string{"react-native", "list-devices", "--json"}
	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.Output()
	if err != nil {
		return HealthCheck{
			Name:    "Devices",
			Status:  StatusWarning,
			Message: "Could not detect connected devices",
			Suggest: "Connect a device via MobAI or USB, then run 'npx react-native list-devices'",
		}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	deviceCount := 0
	for _, line := range lines {
		if strings.Contains(line, "device") || strings.Contains(line, "Device") {
			deviceCount++
		}
	}

	if deviceCount == 0 {
		return HealthCheck{
			Name:    "Devices",
			Status:  StatusWarning,
			Message: "No devices found",
			Suggest: "Connect a device via MobAI or start an iOS simulator",
		}
	}

	return HealthCheck{
		Name:    "Devices",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("%d device(s) available", deviceCount),
	}
}
