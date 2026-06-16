package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FailureCategory string

const (
	CatBuild         FailureCategory = "BUILD"
	CatWorkflow      FailureCategory = "WORKFLOW"
	CatSigning       FailureCategory = "SIGNING"
	CatProvisioning  FailureCategory = "PROVISIONING"
	CatCertificate   FailureCategory = "CERTIFICATE"
	CatGitHubActions FailureCategory = "GITHUB_ACTIONS"
	CatFlutter       FailureCategory = "FLUTTER"
	CatReactNative   FailureCategory = "REACT_NATIVE"
	CatMetro         FailureCategory = "METRO"
	CatDependency    FailureCategory = "DEPENDENCY"
	CatNetwork       FailureCategory = "NETWORK"
	CatPermission    FailureCategory = "PERMISSION"
	CatUnknown       FailureCategory = "UNKNOWN"
)

type AnalysisResult struct {
	Category   FailureCategory `json:"category"`
	Problem    string          `json:"problem"`
	Cause      string          `json:"cause"`
	Impact     string          `json:"impact"`
	Fix        string          `json:"fix"`
	Confidence float64         `json:"confidence"`
	RawLogs    []string        `json:"raw_logs,omitempty"`
}

type Analyzer struct {
	classifiers []Classifier
}

type Classifier interface {
	Name() string
	Match(line string) *FailureCategory
}

type patternClassifier struct {
	name     string
	patterns []string
	category FailureCategory
}

func (p *patternClassifier) Name() string { return p.name }

func (p *patternClassifier) Match(line string) *FailureCategory {
	lower := strings.ToLower(line)
	for _, pat := range p.patterns {
		if strings.Contains(lower, pat) {
			return &p.category
		}
	}
	return nil
}

func NewAnalyzer() *Analyzer {
	classifiers := []Classifier{
		&patternClassifier{
			name:     "signing error",
			patterns: []string{"code signing", "code signing", "provisioning profile", "certificate:", "team id"},
			category: CatSigning,
		},
		&patternClassifier{
			name:     "flutter error",
			patterns: []string{"flutter:", "dart:", "pub get", "flutter run", "flutter:"},
			category: CatFlutter,
		},
		&patternClassifier{
			name:     "react native error",
			patterns: []string{"react-native:", "metro bundler", "metro "},
			category: CatReactNative,
		},
		&patternClassifier{
			name:     "network error",
			patterns: []string{"network error", "timeout:", "connection refused", "no route", "dns:"},
			category: CatNetwork,
		},
		&patternClassifier{
			name:     "dependency error",
			patterns: []string{"cannot find module", "module not found", "cannot find", "missing dependency"},
			category: CatDependency,
		},
		&patternClassifier{
			name:     "permission error",
			patterns: []string{"permission denied", "access denied", "unauthorized", "forbidden"},
			category: CatPermission,
		},
		&patternClassifier{
			name:     "github actions error",
			patterns: []string{"github actions", "workflow error", "runner error", "ci/cd"},
			category: CatGitHubActions,
		},
		&patternClassifier{
			name:     "xcode build error",
			patterns: []string{"error:", "error ", "command failed", "build failed", "xcodebuild failed"},
			category: CatBuild,
		},
	}

	return &Analyzer{classifiers: classifiers}
}

func (a *Analyzer) Analyze(logs []string) *AnalysisResult {
	result := &AnalysisResult{
		Category:   CatUnknown,
		Confidence: 0.0,
		RawLogs:    logs,
	}

	if len(logs) == 0 {
		result.Problem = "No logs available"
		result.Cause = "No build or workflow logs found to analyze"
		result.Impact = "Cannot determine failure cause"
		result.Fix = "Run a build first to generate logs"
		result.Confidence = 0.0
		return result
	}

	categoryScores := make(map[FailureCategory]int)
	var errorLines []string

	for _, log := range logs {
		for _, line := range strings.Split(log, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			for _, c := range a.classifiers {
				if match := c.Match(line); match != nil {
					categoryScores[*match]++
					errorLines = append(errorLines, line)
				}
			}
		}
	}

	if len(categoryScores) == 0 {
		result.Problem = "Build failed (unknown cause)"
		result.Cause = "Could not automatically identify the failure from logs"
		result.Impact = "Build cannot proceed until the issue is resolved"
		result.Fix = "Review the build logs manually or run 'builder ai doctor' for a full audit"
		result.Confidence = 0.1
		return result
	}

	bestCategory := CatUnknown
	bestScore := 0
	total := 0
	classifierOrder := []FailureCategory{
		CatSigning, CatProvisioning, CatCertificate,
		CatFlutter, CatReactNative, CatMetro,
		CatNetwork, CatDependency, CatPermission,
		CatGitHubActions, CatBuild,
	}
	for _, score := range categoryScores {
		total += score
	}
	for _, cat := range classifierOrder {
		if score, ok := categoryScores[cat]; ok && score > bestScore {
			bestScore = score
			bestCategory = cat
		}
	}

	result.Category = bestCategory
	result.Confidence = float64(bestScore) / float64(total)
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}

	result.Problem = fmt.Sprintf("%s failure detected", bestCategory)
	result.Cause = fmt.Sprintf("Found %d error pattern(s) matching %s category from %d log lines",
		bestScore, bestCategory, len(errorLines))
	result.Impact = fmt.Sprintf("The %s issue is blocking the build process", bestCategory)

	kb := NewKnowledgeBase()
	fix := kb.Lookup(bestCategory)
	if fix != "" {
		result.Fix = fix
	} else {
		result.Fix = fmt.Sprintf("Check the %s configuration and logs for details", bestCategory)
	}

	return result
}

type AnalysisEngine struct {
	Analyzer *Analyzer
	KB       *KnowledgeBase
}

func NewAnalysisEngine() *AnalysisEngine {
	return &AnalysisEngine{
		Analyzer: NewAnalyzer(),
		KB:       NewKnowledgeBase(),
	}
}

func (e *AnalysisEngine) AnalyzeBuild(dir string) (*AnalysisResult, error) {
	logsDir := filepath.Join(dir, ".build", "logs")
	var allLogs []string

	err := filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err == nil {
				allLogs = append(allLogs, fmt.Sprintf("[%s] %s", path, string(data)))
			}
		}
		return nil
	})
	if err != nil {
		allLogs = []string{}
	}

	return e.Analyzer.Analyze(allLogs), nil
}
