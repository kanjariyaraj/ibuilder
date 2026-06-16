package reactnative

import (
	"fmt"
	"os"
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

	report.Checks = append(report.Checks, s.checkNodeJS())
	report.Checks = append(report.Checks, s.checkPackageManager())
	report.Checks = append(report.Checks, s.checkRNProject())
	report.Checks = append(report.Checks, s.checkDependencies())
	report.Checks = append(report.Checks, s.checkMetroCLI())
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

func (s *Session) checkNodeJS() HealthCheck {
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
		Message: fmt.Sprintf("Node.js %s installed", strings.TrimSpace(string(output))),
	}
}

func (s *Session) checkPackageManager() HealthCheck {
	hasNPM := true
	cmd := exec.Command("npm", "--version")
	if err := cmd.Run(); err != nil {
		hasNPM = false
	}

	if !hasNPM {
		return HealthCheck{
			Name:    "Package Manager",
			Status:  StatusFailure,
			Message: "npm is not installed",
			Suggest: "npm is bundled with Node.js — reinstall Node.js",
		}
	}

	npmVer, _ := s.CheckNPM()
	var pmDetails string
	yarnVer, _ := s.CheckYarn()
	pnpmVer, _ := s.CheckPNPM()

	pmDetails = fmt.Sprintf("npm %s", npmVer)
	if yarnVer != "" {
		pmDetails += fmt.Sprintf(", yarn %s", yarnVer)
	}
	if pnpmVer != "" {
		pmDetails += fmt.Sprintf(", pnpm %s", pnpmVer)
	}

	return HealthCheck{
		Name:    "Package Manager",
		Status:  StatusHealthy,
		Message: pmDetails,
	}
}

func (s *Session) checkRNProject() HealthCheck {
	if s.projectDir == "" {
		return HealthCheck{
			Name:    "React Native Project",
			Status:  StatusWarning,
			Message: "No project directory set",
			Suggest: "Run 'builder rn dev' from within a React Native project directory",
		}
	}

	valid, err := s.DetectRNProject(s.projectDir)
	if err != nil {
		return HealthCheck{
			Name:    "React Native Project",
			Status:  StatusFailure,
			Message: fmt.Sprintf("Invalid project: %v", err),
			Suggest: "Ensure the directory contains a valid React Native iOS project with package.json (containing react-native dependency) and ios/",
		}
	}
	if !valid {
		return HealthCheck{
			Name:    "React Native Project",
			Status:  StatusFailure,
			Message: "Not a valid React Native project",
			Suggest: "Run 'npx react-native init' or navigate to a React Native project directory",
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
	if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
		return HealthCheck{
			Name:    "Dependencies",
			Status:  StatusWarning,
			Message: "node_modules not found — dependencies not installed",
			Suggest: "Run 'npm install' or 'yarn install' in your project directory",
		}
	}

	return HealthCheck{
		Name:    "Dependencies",
		Status:  StatusHealthy,
		Message: "node_modules present",
	}
}

func (s *Session) checkMetroCLI() HealthCheck {
	if !s.CheckMetro() {
		return HealthCheck{
			Name:    "Metro Bundler",
			Status:  StatusWarning,
			Message: "Metro CLI not detected",
			Suggest: "Ensure react-native is installed in your project dependencies",
		}
	}

	return HealthCheck{
		Name:    "Metro Bundler",
		Status:  StatusHealthy,
		Message: "Metro CLI available",
	}
}

func (s *Session) checkDevices() HealthCheck {
	cmd := exec.Command("npx", "react-native", "list-devices")
	output, err := cmd.Output()
	if err != nil {
		return HealthCheck{
			Name:    "Devices",
			Status:  StatusWarning,
			Message: "Could not detect connected devices",
			Suggest: "Connect a device via MobAI or USB, then check with 'npx react-native list-devices'",
		}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	deviceCount := 0
	for _, line := range lines {
		if strings.Contains(line, "device") || strings.Contains(line, "emulator") {
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
