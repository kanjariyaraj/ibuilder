package flutter

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

	report.Checks = append(report.Checks, s.checkFlutterSDK())
	report.Checks = append(report.Checks, s.checkDartSDK())
	report.Checks = append(report.Checks, s.checkProject())
	report.Checks = append(report.Checks, s.checkDependencies())
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

func (s *Session) checkFlutterSDK() HealthCheck {
	cmd := exec.Command("flutter", "--version")
	output, err := cmd.Output()
	if err != nil {
		return HealthCheck{
			Name:    "Flutter SDK",
			Status:  StatusFailure,
			Message: "Flutter is not installed or not in PATH",
			Suggest: "Install Flutter from https://flutter.dev or run 'flutter doctor'",
		}
	}

	lines := strings.Split(string(output), "\n")
	version := "unknown"
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) > 1 {
			version = parts[1]
		}
	}

	return HealthCheck{
		Name:    "Flutter SDK",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("Flutter %s installed", version),
	}
}

func (s *Session) checkDartSDK() HealthCheck {
	cmd := exec.Command("dart", "--version")
	output, err := cmd.Output()
	if err != nil {
		return HealthCheck{
			Name:    "Dart SDK",
			Status:  StatusFailure,
			Message: "Dart SDK is not installed or not in PATH",
			Suggest: "Dart is bundled with Flutter — ensure Flutter is properly installed",
		}
	}

	version := strings.TrimSpace(string(output))
	return HealthCheck{
		Name:    "Dart SDK",
		Status:  StatusHealthy,
		Message: version,
	}
}

func (s *Session) checkProject() HealthCheck {
	if s.projectDir == "" {
		return HealthCheck{
			Name:    "Flutter Project",
			Status:  StatusWarning,
			Message: "No project directory set",
			Suggest: "Run 'builder flutter dev' from within a Flutter project directory",
		}
	}

	valid, err := s.DetectFlutterProject(s.projectDir)
	if err != nil {
		return HealthCheck{
			Name:    "Flutter Project",
			Status:  StatusFailure,
			Message: fmt.Sprintf("Invalid project: %v", err),
			Suggest: "Ensure the directory contains a valid Flutter iOS project with pubspec.yaml and ios/",
		}
	}
	if !valid {
		return HealthCheck{
			Name:    "Flutter Project",
			Status:  StatusFailure,
			Message: "Not a valid Flutter project",
			Suggest: "Run 'flutter create' or navigate to a Flutter project directory",
		}
	}

	return HealthCheck{
		Name:    "Flutter Project",
		Status:  StatusHealthy,
		Message: fmt.Sprintf("Valid Flutter project in %s", s.projectDir),
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

	cmd := exec.Command("flutter", "pub", "check")
	cmd.Dir = s.projectDir
	if err := cmd.Run(); err != nil {
		return HealthCheck{
			Name:    "Dependencies",
			Status:  StatusWarning,
			Message: "Dependencies may need resolution",
			Suggest: "Run 'flutter pub get' in your project directory",
		}
	}

	return HealthCheck{
		Name:    "Dependencies",
		Status:  StatusHealthy,
		Message: "Dependencies resolved",
	}
}

func (s *Session) checkDevices() HealthCheck {
	cmd := exec.Command("flutter", "devices")
	output, err := cmd.Output()
	if err != nil {
		return HealthCheck{
			Name:    "Devices",
			Status:  StatusWarning,
			Message: "Could not detect Flutter devices",
			Suggest: "Connect a device or start an emulator, then run 'flutter devices'",
		}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	deviceCount := 0
	for _, line := range lines {
		if strings.Contains(line, "•") || strings.Contains(line, "*") {
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
