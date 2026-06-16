package init

import (
	"os"
	"path/filepath"
	"strings"
)

func DetectSchemesAndWorkspaces(dir string) (schemes []string, workspaces []string, err error) {
	walkDirs := []string{dir}
	iosDir := filepath.Join(dir, "ios")
	if hasDir(dir, "ios") {
		walkDirs = append(walkDirs, iosDir)
	}

	for _, d := range walkDirs {
		ws, _ := filepath.Glob(filepath.Join(d, "*.xcworkspace"))
		for _, w := range ws {
			workspaces = append(workspaces, strings.TrimSuffix(filepath.Base(w), ".xcworkspace"))
		}

		xp, _ := filepath.Glob(filepath.Join(d, "*.xcodeproj"))
		for _, x := range xp {
			pbxPath := filepath.Join(x, "project.pbxproj")
			if data, err := os.ReadFile(pbxPath); err == nil {
				s := extractTargets(string(data))
				schemes = append(schemes, s...)
			}
		}
	}

	schemes = unique(schemes)
	return
}

func extractTargets(pbxContent string) []string {
	var targets []string
	lines := strings.Split(pbxContent, "\n")
	inTargets := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "/* Begin PBXNativeTarget section */") {
			inTargets = true
			continue
		}
		if strings.Contains(line, "/* End PBXNativeTarget section */") {
			break
		}
		if inTargets && strings.Contains(line, "isa = PBXNativeTarget") {
			parts := strings.Split(line, "/*")
			if len(parts) >= 2 {
				name := strings.TrimSpace(strings.Split(parts[1], "*/")[0])
				if name != "" {
					targets = append(targets, name)
				}
			}
		}
	}
	return targets
}

func unique(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func (s *schemeDetector) DetectSchemes(dir string) ([]string, error) {
	schemes, _, err := DetectSchemesAndWorkspaces(dir)
	return schemes, err
}

func (s *schemeDetector) DetectWorkspaces(dir string) ([]string, error) {
	_, workspaces, err := DetectSchemesAndWorkspaces(dir)
	return workspaces, err
}

type schemeDetector struct{}
