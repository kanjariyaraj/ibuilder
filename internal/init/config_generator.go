package init

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kanjariyaraj/Builder/internal/config"
)

func generateConfig(pt ProjectType, projectName, repoOwner, repoName, dir string) *config.Config {
	cfg := config.Default()
	cfg.ProjectName = projectName
	cfg.Repo.Owner = repoOwner
	cfg.Repo.Name = repoName

	iosPath := detectIOSPath(dir, pt)
	workspace := detectWorkspace(dir)
	schemes := detectSchemes(dir)

	cfg.Repository = ""
	if repoOwner != "" && repoName != "" {
		cfg.Repository = "https://github.com/" + repoOwner + "/" + repoName + ".git"
	}

	switch pt {
	case ProjectFlutter:
		cfg.Flutter.Enabled = true
		cfg.Build.ProjectType = "flutter"
	case ProjectReactNative:
		cfg.ReactNative.Enabled = true
		cfg.Build.ProjectType = "react-native"
	case ProjectExpo:
		cfg.ReactNative.Enabled = true
		cfg.Build.ProjectType = "expo"
	case ProjectCapacitor:
		cfg.Build.ProjectType = "capacitor"
	case ProjectNativeiOS:
		cfg.Build.ProjectType = "xcode"
	}

	if iosPath != "" {
		cfg.IOS.MinimumVersion = "15.0"
		cfg.IOS.TargetVersion = "17.0"
	}
	if workspace != "" {
		cfg.Build.Scheme = workspace
	}
	if len(schemes) > 0 {
		if cfg.Build.Scheme == "" {
			cfg.Build.Scheme = schemes[0]
		}
	}

	return cfg
}

func detectIOSPath(dir string, pt ProjectType) string {
	switch pt {
	case ProjectFlutter:
		if hasDir(dir, "ios") {
			return "ios"
		}
	case ProjectReactNative, ProjectExpo:
		if hasDir(dir, "ios") {
			return "ios"
		}
		return ""
	case ProjectCapacitor:
		if hasDir(filepath.Join(dir, "ios"), "App") {
			return "ios"
		}
	case ProjectNativeiOS:
		return "."
	}
	return ""
}

func detectWorkspace(dir string) string {
	iosDir := filepath.Join(dir, "ios")
	matches, _ := filepath.Glob(filepath.Join(dir, "*.xcworkspace"))
	if len(matches) > 0 {
		return strings.TrimSuffix(filepath.Base(matches[0]), ".xcworkspace")
	}
	matches, _ = filepath.Glob(filepath.Join(iosDir, "*.xcworkspace"))
	if len(matches) > 0 {
		return strings.TrimSuffix(filepath.Base(matches[0]), ".xcworkspace")
	}
	return ""
}

func detectSchemes(dir string) []string {
	var schemes []string
	content, err := os.ReadFile(filepath.Join(dir, "Runner.xcodeproj", "project.pbxproj"))
	if err == nil {
		sc := string(content)
		if strings.Contains(sc, "Runner") {
			schemes = append(schemes, "Runner")
		}
	}
	content, err = os.ReadFile(filepath.Join(dir, "*.xcodeproj", "project.pbxproj"))
	if err == nil {
		sc := string(content)
		if strings.Contains(sc, "App") {
			schemes = append(schemes, "App")
		}
	}
	return schemes
}
