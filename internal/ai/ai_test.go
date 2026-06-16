package ai

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	s := NewSession(nil)
	if s == nil {
		t.Fatal("expected session, got nil")
	}
}

func TestNewAnalyzer(t *testing.T) {
	a := NewAnalyzer()
	if a == nil {
		t.Fatal("expected analyzer, got nil")
	}
}

func TestAnalyzeNoLogs(t *testing.T) {
	a := NewAnalyzer()
	result := a.Analyze(nil)
	if result.Confidence != 0.0 {
		t.Errorf("expected confidence 0.0, got %f", result.Confidence)
	}
}

func TestAnalyzeEmptyLogs(t *testing.T) {
	a := NewAnalyzer()
	result := a.Analyze([]string{})
	if result.Category != CatUnknown {
		t.Errorf("expected CatUnknown, got %s", result.Category)
	}
}

func TestAnalyzeBuildError(t *testing.T) {
	a := NewAnalyzer()
	logs := []string{"xcodebuild failed with error: Command failed"}
	result := a.Analyze(logs)
	if result.Category != CatBuild {
		t.Errorf("expected CatBuild, got %s", result.Category)
	}
}

func TestAnalyzeSigningError(t *testing.T) {
	a := NewAnalyzer()
	logs := []string{"Code signing error: No matching provisioning profile found"}
	result := a.Analyze(logs)
	if result.Category != CatSigning {
		t.Errorf("expected CatSigning, got %s", result.Category)
	}
}

func TestAnalyzeFlutterError(t *testing.T) {
	a := NewAnalyzer()
	logs := []string{"flutter: The method 'foo' isn't defined"}
	result := a.Analyze(logs)
	if result.Category != CatFlutter {
		t.Errorf("expected CatFlutter, got %s", result.Category)
	}
}

func TestAnalyzeReactNativeError(t *testing.T) {
	a := NewAnalyzer()
	logs := []string{"react-native: Metro bundler encountered an error"}
	result := a.Analyze(logs)
	if result.Category != CatReactNative {
		t.Errorf("expected CatReactNative, got %s", result.Category)
	}
}

func TestAnalyzeNetworkError(t *testing.T) {
	a := NewAnalyzer()
	logs := []string{"network error: connection refused to github.com"}
	result := a.Analyze(logs)
	if result.Category != CatNetwork {
		t.Errorf("expected CatNetwork, got %s", result.Category)
	}
}

func TestAnalyzeDependencyError(t *testing.T) {
	a := NewAnalyzer()
	logs := []string{"Error: Cannot find module 'react-native'"}
	result := a.Analyze(logs)
	if result.Category != CatDependency {
		t.Errorf("expected CatDependency, got %s", result.Category)
	}
}

func TestAnalyzeMultipleLines(t *testing.T) {
	a := NewAnalyzer()
	logs := []string{
		"Started build at ...",
		"Error: xcodebuild failed with code 65",
		"Code signing is required for product type Application",
	}
	result := a.Analyze(logs)
	if result.Confidence <= 0 {
		t.Errorf("expected positive confidence, got %f", result.Confidence)
	}
}

func TestKnowledgeBaseLookup(t *testing.T) {
	kb := NewKnowledgeBase()
	fix := kb.Lookup(CatBuild)
	if fix == "" {
		t.Error("expected non-empty fix for CatBuild")
	}
}

func TestKnowledgeBaseUnknown(t *testing.T) {
	kb := NewKnowledgeBase()
	fix := kb.Lookup("UNKNOWN_CATEGORY")
	if fix != "" {
		t.Errorf("expected empty fix for unknown category, got %s", fix)
	}
}

func TestKnowledgeBaseAddRule(t *testing.T) {
	kb := NewKnowledgeBase()
	kb.AddRule(CatBuild, "Extra build rule")
	fix := kb.Lookup(CatBuild)
	if fix == "" {
		t.Error("expected non-empty fix after adding rule")
	}
}

func TestKnowledgeBaseCategories(t *testing.T) {
	kb := NewKnowledgeBase()
	cats := kb.Categories()
	if len(cats) == 0 {
		t.Error("expected non-empty categories")
	}
}

func TestKnowledgeBaseSearch(t *testing.T) {
	kb := NewKnowledgeBase()
	results := kb.Search("Xcode")
	if len(results) == 0 {
		t.Error("expected search results for 'Xcode'")
	}
}

func TestClassificationReport(t *testing.T) {
	a := NewAnalyzer()
	logs := []string{"Error: build failed", "normal line", "Error: signing issue"}
	reports := a.ClassifyLogs(logs)
	if len(reports) == 0 {
		t.Error("expected at least one classification report")
	}
}

func TestPatternClassifier(t *testing.T) {
	c := &patternClassifier{
		name:     "test",
		patterns: []string{"error", "fail"},
		category: CatBuild,
	}

	if match := c.Match("this is an error"); match == nil || *match != CatBuild {
		t.Error("expected match for 'error'")
	}
	if match := c.Match("everything is fine"); match != nil {
		t.Error("expected no match for normal line")
	}
}

func TestDoctorReport(t *testing.T) {
	s := NewSession(nil)
	report := s.Doctor()
	if len(report.Checks) == 0 {
		t.Error("expected non-empty doctor checks")
	}
}

func TestNewAnalysisEngine(t *testing.T) {
	e := NewAnalysisEngine()
	if e == nil {
		t.Fatal("expected engine, got nil")
	}
	if e.Analyzer == nil {
		t.Error("expected analyzer in engine")
	}
	if e.KB == nil {
		t.Error("expected knowledge base in engine")
	}
}

func TestFailureCategoryConstants(t *testing.T) {
	if CatBuild != "BUILD" {
		t.Errorf("expected BUILD, got %s", CatBuild)
	}
	if CatSigning != "SIGNING" {
		t.Errorf("expected SIGNING, got %s", CatSigning)
	}
	if CatUnknown != "UNKNOWN" {
		t.Errorf("expected UNKNOWN, got %s", CatUnknown)
	}
}

func TestReportFormatConstants(t *testing.T) {
	if FormatMarkdown != "markdown" {
		t.Errorf("expected markdown, got %s", FormatMarkdown)
	}
	if FormatJSON != "json" {
		t.Errorf("expected json, got %s", FormatJSON)
	}
	if FormatHTML != "html" {
		t.Errorf("expected html, got %s", FormatHTML)
	}
}
