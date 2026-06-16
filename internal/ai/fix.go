package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FixSuggestion struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Current string `json:"current"`
	Suggest string `json:"suggest"`
	Reason  string `json:"reason"`
}

type AutoFixResult struct {
	Success     bool            `json:"success"`
	Action      string          `json:"action"`
	Fixes       []FixSuggestion `json:"fixes,omitempty"`
	DryRun      bool            `json:"dry_run"`
	Summary     string          `json:"summary,omitempty"`
	Error       string          `json:"error,omitempty"`
}

func (s *Session) AnalyzeAndSuggestFixes(dir string) (*AutoFixResult, error) {
	s.log.Info("analyzing for auto-fix suggestions", "dir", dir)

	logs, err := s.collectLogs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to collect logs: %w", err)
	}

	analysis := s.analyzer.Analyze(logs)

	var fixes []FixSuggestion
	switch analysis.Category {
	case CatBuild:
		fixes = s.suggestBuildFixes(dir)
	case CatSigning:
		fixes = s.suggestSigningFixes(dir)
	case CatDependency:
		fixes = s.suggestDependencyFixes(dir)
	case CatFlutter:
		fixes = s.suggestFlutterFixes(dir)
	case CatReactNative:
		fixes = s.suggestReactNativeFixes(dir)
	case CatNetwork:
		fixes = s.suggestNetworkFixes(dir)
	}

	result := &AutoFixResult{
		Success: true,
		Action:  "analysis",
		Fixes:   fixes,
		DryRun:  true,
	}

	if len(fixes) == 0 {
		result.Summary = fmt.Sprintf("Analysis complete: %s. No automatic fixes available.", analysis.Problem)
		result.Success = false
	} else {
		result.Summary = fmt.Sprintf("Found %d potential fix(es) for %s. Review with 'builder ai fix --apply'.", len(fixes), analysis.Category)
	}

	return result, nil
}

func (s *Session) ApplyFix(dir string, fix FixSuggestion, dryRun bool) (*AutoFixResult, error) {
	s.log.Info("applying fix", "file", fix.File, "dry_run", dryRun)

	fullPath := filepath.Join(dir, fix.File)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", fix.File, err)
	}

	content := string(data)
	if dryRun {
		result := &AutoFixResult{
			Success: true,
			Action:  "dry_run",
			Fixes:   []FixSuggestion{fix},
			DryRun:  true,
			Summary: fmt.Sprintf("[DRY RUN] Would replace '%s' with '%s' in %s",
				truncate(fix.Current, 40), truncate(fix.Suggest, 40), fix.File),
		}
		return result, nil
	}

	if !strings.Contains(content, fix.Current) {
		return nil, fmt.Errorf("could not find the expected content in %s", fix.File)
	}

	newContent := strings.Replace(content, fix.Current, fix.Suggest, 1)
	if err := os.WriteFile(fullPath, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write %s: %w", fix.File, err)
	}

	result := &AutoFixResult{
		Success: true,
		Action:  "applied",
		Fixes:   []FixSuggestion{fix},
		DryRun:  false,
		Summary: fmt.Sprintf("Applied fix to %s", fix.File),
	}

	s.log.Info("fix applied", "file", fix.File)
	return result, nil
}

func (s *Session) suggestBuildFixes(dir string) []FixSuggestion {
	cfgPath := filepath.Join(dir, "builder.json")
	if _, err := os.Stat(cfgPath); err != nil {
		return nil
	}
	data, _ := os.ReadFile(cfgPath)
	content := string(data)

	var fixes []FixSuggestion
	if strings.Contains(content, `"build_mode": "debug"`) {
		fixes = append(fixes, FixSuggestion{
			File:    "builder.json",
			Current: `"build_mode": "debug"`,
			Suggest: `"build_mode": "release"`,
			Reason:  "Release mode is recommended for production builds",
		})
	}

	return fixes
}

func (s *Session) suggestSigningFixes(dir string) []FixSuggestion {
	cfgPath := filepath.Join(dir, "builder.json")
	if _, err := os.Stat(cfgPath); err != nil {
		return nil
	}
	data, _ := os.ReadFile(cfgPath)
	content := string(data)

	var fixes []FixSuggestion
	if !strings.Contains(content, `"team_id"`) {
		fixes = append(fixes, FixSuggestion{
			File:    "builder.json",
			Current: `"signing": {`,
			Suggest: `"signing": {\n    "team_id": "<YOUR_TEAM_ID>"`,
			Reason:  "Team ID is required for code signing",
		})
	}

	return fixes
}

func (s *Session) suggestDependencyFixes(dir string) []FixSuggestion {
	pkgPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(pkgPath); err != nil {
		return nil
	}
	data, _ := os.ReadFile(pkgPath)
	content := string(data)

	var fixes []FixSuggestion
	if !strings.Contains(content, `"node_modules"`) {
		fixes = append(fixes, FixSuggestion{
			File:    "package.json",
			Current: "",
			Suggest: "Run 'npm install' to install dependencies",
			Reason:  "Dependencies are not installed",
		})
	}

	return fixes
}

func (s *Session) suggestFlutterFixes(dir string) []FixSuggestion {
	var fixes []FixSuggestion
	cfgPath := filepath.Join(dir, "builder.json")
	data, _ := os.ReadFile(cfgPath)
	content := string(data)

	if !strings.Contains(content, `"flutter"`) {
		fixes = append(fixes, FixSuggestion{
			File:    "builder.json",
			Current: "flutter configuration missing",
			Suggest: "Add flutter section to builder.json",
			Reason:  "Flutter project detected but not configured in builder.json",
		})
	}

	return fixes
}

func (s *Session) suggestReactNativeFixes(dir string) []FixSuggestion {
	var fixes []FixSuggestion
	cfgPath := filepath.Join(dir, "builder.json")
	data, _ := os.ReadFile(cfgPath)
	content := string(data)

	if !strings.Contains(content, `"react_native"`) {
		fixes = append(fixes, FixSuggestion{
			File:    "builder.json",
			Current: "react_native configuration missing",
			Suggest: "Add react_native section to builder.json",
			Reason:  "React Native project detected but not configured",
		})
	}

	return fixes
}

func (s *Session) suggestNetworkFixes(dir string) []FixSuggestion {
	return []FixSuggestion{
		{
			File:    "environment",
			Current: "Network connectivity issue",
			Suggest: "Check internet connection and proxy settings",
			Reason:  "Network failures prevent dependency downloads and API calls",
		},
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
