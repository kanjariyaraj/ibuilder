package init

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFlutter(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "pubspec.yaml"), []byte("name: test_app"), 0644)

	pt := DetectProjectType(dir)
	if pt != ProjectFlutter {
		t.Errorf("expected Flutter, got %s", pt)
	}
}

func TestDetectReactNative(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"dependencies":{"react-native":"0.73"}}`), 0644)

	pt := DetectProjectType(dir)
	if pt != ProjectReactNative {
		t.Errorf("expected ReactNative, got %s", pt)
	}
}

func TestDetectExpo(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"dependencies":{"expo":"50.0"}}`), 0644)

	pt := DetectProjectType(dir)
	if pt != ProjectExpo {
		t.Errorf("expected Expo, got %s", pt)
	}
}

func TestDetectCapacitor(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"dependencies":{"@capacitor/core":"5.0"}}`), 0644)

	pt := DetectProjectType(dir)
	if pt != ProjectCapacitor {
		t.Errorf("expected Capacitor, got %s", pt)
	}
}

func TestDetectNativeiOS(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "MyApp.xcodeproj"), 0755)

	pt := DetectProjectType(dir)
	if pt != ProjectNativeiOS {
		t.Errorf("expected NativeiOS, got %s", pt)
	}
}

func TestDetectUnknown(t *testing.T) {
	dir := t.TempDir()
	pt := DetectProjectType(dir)
	if pt != ProjectUnknown {
		t.Errorf("expected Unknown, got %s", pt)
	}
}

func TestDetectProjectNameFlutter(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "pubspec.yaml"), []byte("name: my_flutter_app\n"), 0644)
	name := detectProjectName(dir, ProjectFlutter)
	if name != "my_flutter_app" {
		t.Errorf("expected 'my_flutter_app', got '%s'", name)
	}
}

func TestDetectProjectNameReactNative(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"MyRNApp"}`), 0644)
	name := detectProjectName(dir, ProjectReactNative)
	if name != "MyRNApp" {
		t.Errorf("expected 'MyRNApp', got '%s'", name)
	}
}

func TestGenerateConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := generateConfig(ProjectFlutter, "test_app", "owner", "repo", dir)
	if cfg.ProjectName != "test_app" {
		t.Errorf("expected 'test_app', got '%s'", cfg.ProjectName)
	}
	if !cfg.Flutter.Enabled {
		t.Errorf("expected Flutter.Enabled true")
	}
	if cfg.Repo.Owner != "owner" {
		t.Errorf("expected 'owner', got '%s'", cfg.Repo.Owner)
	}
}

func TestFlutterTemplate(t *testing.T) {
	tmpl := getTemplate(ProjectFlutter, "MyApp")
	if tmpl == "" {
		t.Errorf("expected non-empty template")
	}
}

func TestNativeIOSTemplate(t *testing.T) {
	tmpl := getTemplate(ProjectNativeiOS, "MyApp")
	if tmpl == "" {
		t.Errorf("expected non-empty template")
	}
}

func TestValidateGeneratedFiles(t *testing.T) {
	dir := t.TempDir()
	errs := ValidateGeneratedFiles(dir)
	if len(errs) == 0 {
		t.Errorf("expected validation errors for empty directory")
	}
}

func TestExtractTargets(t *testing.T) {
	pbx := `/* Begin PBXNativeTarget section */
		1A2B /* MyApp */ = { isa = PBXNativeTarget; };
		3C4D /* WidgetExtension */ = { isa = PBXNativeTarget; };
/* End PBXNativeTarget section */`
	targets := extractTargets(pbx)
	if len(targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(targets))
	}
	if targets[0] != "MyApp" {
		t.Errorf("expected 'MyApp', got '%s'", targets[0])
	}
}

func TestUnique(t *testing.T) {
	items := unique([]string{"a", "b", "a", "c"})
	if len(items) != 3 {
		t.Errorf("expected 3 unique items, got %d", len(items))
	}
}
