package init

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DetectProjectType(dir string) ProjectType {
	if hasFile(dir, "pubspec.yaml") {
		return ProjectFlutter
	}
	if hasFile(dir, "package.json") {
		pkg := readFile(dir, "package.json")
		if strings.Contains(pkg, "\"expo\"") {
			return ProjectExpo
		}
		if strings.Contains(pkg, "\"react-native\"") {
			return ProjectReactNative
		}
		if strings.Contains(pkg, "\"@capacitor") {
			return ProjectCapacitor
		}
		if strings.Contains(pkg, "\"@ionic") {
			return ProjectIonic
		}
		if strings.Contains(pkg, "\"cordova") {
			return ProjectCordova
		}
	}
	if hasGlob(dir, "*.xcworkspace") || hasGlob(dir, "*.xcodeproj") {
		return ProjectNativeiOS
	}
	if hasFile(dir, "Assembly-CSharp.csproj") || hasDir(dir, "Assets") && hasDir(dir, "ProjectSettings") {
		return ProjectUnity
	}
	if hasFile(dir, "*.uproject") || hasDir(dir, "Source") && hasDir(dir, "Config") {
		return ProjectUnreal
	}
	if hasDir(dir, "shared") && hasDir(dir, "androidApp") && hasDir(dir, "iosApp") {
		return ProjectKMM
	}
	return ProjectUnknown
}

func detectProjectName(dir string, pt ProjectType) string {
	switch pt {
	case ProjectFlutter:
		return readNameFromYAML(dir, "pubspec.yaml")
	case ProjectReactNative, ProjectExpo, ProjectCapacitor, ProjectIonic, ProjectCordova:
		return readNameFromJSON(dir, "package.json")
	case ProjectNativeiOS:
		return detectXcodeProjectName(dir)
	case ProjectUnity:
		return filepath.Base(dir)
	case ProjectUnreal:
		return filepath.Base(dir)
	case ProjectKMM:
		return filepath.Base(dir)
	default:
		return filepath.Base(dir)
	}
}

func detectRepoInfo(dir string) (string, string) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", ""
	}
	url := strings.TrimSpace(string(out))
	if strings.HasPrefix(url, "https://") {
		parts := strings.Split(strings.TrimPrefix(url, "https://"), "/")
		if len(parts) >= 2 {
			return parts[len(parts)-2], strings.TrimSuffix(parts[len(parts)-1], ".git")
		}
	}
	if strings.HasPrefix(url, "git@") {
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			return parts[len(parts)-2], strings.TrimSuffix(parts[len(parts)-1], ".git")
		}
	}
	return "", ""
}

func hasFile(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

func hasDir(dir, name string) bool {
	info, err := os.Stat(filepath.Join(dir, name))
	return err == nil && info.IsDir()
}

func hasGlob(dir, pattern string) bool {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	return err == nil && len(matches) > 0
}

func readFile(dir, name string) string {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return ""
	}
	return string(data)
}

func readNameFromJSON(dir, file string) string {
	content := readFile(dir, file)
	idx := strings.Index(content, "\"name\"")
	if idx < 0 {
		return filepath.Base(dir)
	}
	rest := content[idx+len("\"name\""):]
	colonIdx := strings.Index(rest, ":")
	if colonIdx < 0 {
		return filepath.Base(dir)
	}
	valStart := colonIdx + 1
	val := strings.TrimSpace(rest[valStart:])
	val = strings.Trim(val, "\",}")
	val = strings.Trim(val, "\"")
	if val == "" {
		return filepath.Base(dir)
	}
	return val
}

func readNameFromYAML(dir, file string) string {
	content := readFile(dir, file)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		}
	}
	return filepath.Base(dir)
}

func detectXcodeProjectName(dir string) string {
	matches, _ := filepath.Glob(filepath.Join(dir, "*.xcodeproj"))
	if len(matches) > 0 {
		return strings.TrimSuffix(filepath.Base(matches[0]), ".xcodeproj")
	}
	matches, _ = filepath.Glob(filepath.Join(dir, "*.xcworkspace"))
	if len(matches) > 0 {
		return strings.TrimSuffix(filepath.Base(matches[0]), ".xcworkspace")
	}
	return filepath.Base(dir)
}
