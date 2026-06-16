package ai

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DoctorCheck struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Suggest string `json:"suggest,omitempty"`
}

type DoctorReport struct {
	Healthy bool          `json:"healthy"`
	Checks  []DoctorCheck `json:"checks"`
}

func (s *Session) Doctor() DoctorReport {
	report := DoctorReport{Checks: []DoctorCheck{}}

	report.Checks = append(report.Checks, checkNode())
	report.Checks = append(report.Checks, checkGo())
	report.Checks = append(report.Checks, checkGit())
	report.Checks = append(report.Checks, s.checkProjectConfig())
	report.Checks = append(report.Checks, s.checkBuildLogs())
	report.Checks = append(report.Checks, s.checkDependencyHealth())

	for _, ch := range report.Checks {
		if ch.Status == "FAILURE" {
			report.Healthy = false
			return report
		}
	}
	report.Healthy = true
	return report
}

func checkNode() DoctorCheck {
	cmd := exec.Command("node", "--version")
	if err := cmd.Run(); err != nil {
		return DoctorCheck{
			Name:    "Node.js",
			Status:  "WARNING",
			Message: "Node.js not found in PATH",
			Suggest: "Required for React Native projects — install from https://nodejs.org",
		}
	}
	output, _ := cmd.Output()
	return DoctorCheck{
		Name:    "Node.js",
		Status:  "HEALTHY",
		Message: fmt.Sprintf("Node.js %s", strings.TrimSpace(string(output))),
	}
}

func checkGo() DoctorCheck {
	cmd := exec.Command("go", "version")
	if err := cmd.Run(); err != nil {
		return DoctorCheck{
			Name:    "Go",
			Status:  "FAILURE",
			Message: "Go not found in PATH",
			Suggest: "Install Go from https://go.dev",
		}
	}
	output, _ := cmd.Output()
	return DoctorCheck{
		Name:    "Go",
		Status:  "HEALTHY",
		Message: fmt.Sprintf("%s", strings.TrimSpace(string(output))),
	}
}

func checkGit() DoctorCheck {
	cmd := exec.Command("git", "version")
	if err := cmd.Run(); err != nil {
		return DoctorCheck{
			Name:    "Git",
			Status:  "FAILURE",
			Message: "Git not found in PATH",
			Suggest: "Install Git from https://git-scm.com",
		}
	}
	output, _ := cmd.Output()
	return DoctorCheck{
		Name:    "Git",
		Status:  "HEALTHY",
		Message: strings.TrimSpace(string(output)),
	}
}

func (s *Session) checkProjectConfig() DoctorCheck {
	dir := s.ProjectDir()
	if dir == "" {
		return DoctorCheck{
			Name:    "Project Config",
			Status:  "WARNING",
			Message: "No project directory set",
			Suggest: "Run from within a Builder project directory",
		}
	}

	configPath := filepath.Join(dir, "builder.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DoctorCheck{
			Name:    "Project Config",
			Status:  "WARNING",
			Message: "builder.json not found",
			Suggest: "Run 'builder init' to create a project configuration",
		}
	}

	return DoctorCheck{
		Name:    "Project Config",
		Status:  "HEALTHY",
		Message: fmt.Sprintf("builder.json found at %s", configPath),
	}
}

func (s *Session) checkBuildLogs() DoctorCheck {
	dir := s.ProjectDir()
	if dir == "" {
		return DoctorCheck{
			Name:    "Build Logs",
			Status:  "WARNING",
			Message: "No project directory set",
		}
	}

	logsDir := filepath.Join(dir, ".build", "logs")
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return DoctorCheck{
			Name:    "Build Logs",
			Status:  "WARNING",
			Message: "No build logs found",
			Suggest: "Run a build first to generate logs for analysis",
		}
	}

	count := len(entries)
	if count > 0 {
		analyzer := NewAnalyzer()
		allLogs := []string{}
		for _, e := range entries {
			if !e.IsDir() {
				data, _ := os.ReadFile(filepath.Join(logsDir, e.Name()))
				allLogs = append(allLogs, string(data))
			}
		}
		result := analyzer.Analyze(allLogs)
		if result.Category != CatUnknown {
			return DoctorCheck{
				Name:    "Build Logs",
				Status:  "WARNING",
				Message: fmt.Sprintf("Found %d log file(s) with %s issues", count, result.Category),
				Suggest: fmt.Sprintf("Run 'builder ai explain' for details. Suggested fix: %s", result.Fix),
			}
		}
	}

	return DoctorCheck{
		Name:    "Build Logs",
		Status:  "HEALTHY",
		Message: fmt.Sprintf("%d log file(s) available, no critical issues detected", count),
	}
}

func (s *Session) checkDependencyHealth() DoctorCheck {
	dir := s.ProjectDir()
	if dir == "" {
		return DoctorCheck{
			Name:    "Dependencies",
			Status:  "WARNING",
			Message: "No project directory set",
		}
	}

	nm := filepath.Join(dir, "node_modules")
	if _, err := os.Stat(nm); err == nil {
		return DoctorCheck{
			Name:    "Dependencies",
			Status:  "HEALTHY",
			Message: "node_modules found",
		}
	}

	goMod := filepath.Join(dir, "go.sum")
	if _, err := os.Stat(goMod); err == nil {
		return DoctorCheck{
			Name:    "Dependencies",
			Status:  "HEALTHY",
			Message: "Go modules found",
		}
	}

	return DoctorCheck{
		Name:    "Dependencies",
		Status:  "WARNING",
		Message: "No dependency artifacts detected",
		Suggest: "Run 'npm install' or 'go mod download' as applicable",
	}
}
