package init

import (
	"fmt"
	"os"
	"strings"

	"github.com/kanjariyaraj/Builder/internal/config"
)

type ValidationError struct {
	Field   string
	Message string
}

func ValidateInit(dir string, cfg *config.Config) []ValidationError {
	var errs []ValidationError

	if !hasDir(dir, ".git") && !hasFile(dir, ".git") {
		errs = append(errs, ValidationError{"git", "not a git repository"})
	}

	if cfg.ProjectName == "" {
		errs = append(errs, ValidationError{"project_name", "project name is required"})
	}

	if cfg.Repo.Owner == "" {
		errs = append(errs, ValidationError{"repo.owner", "repository owner not detected"})
	}

	if cfg.Build.ProjectType == "" {
		errs = append(errs, ValidationError{"build.project_type", "project type not set"})
	}

	return errs
}

func ValidateGeneratedFiles(dir string) []ValidationError {
	var errs []ValidationError

	cfgPath := dir + "/builder.json"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		errs = append(errs, ValidationError{"config", "builder.json not generated"})
	}

	wfPath := dir + "/.github/workflows/ios-build.yml"
	if _, err := os.Stat(wfPath); os.IsNotExist(err) {
		errs = append(errs, ValidationError{"workflow", ".github/workflows/ios-build.yml not generated"})
	}

	if _, err := os.Stat(wfPath); err == nil {
		data, _ := os.ReadFile(wfPath)
		content := string(data)
		if !strings.Contains(content, "on:") {
			errs = append(errs, ValidationError{"workflow", "invalid workflow YAML: missing trigger"})
		}
		if !strings.Contains(content, "jobs:") {
			errs = append(errs, ValidationError{"workflow", "invalid workflow YAML: missing jobs"})
		}
	}

	return errs
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
