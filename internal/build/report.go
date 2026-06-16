package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func generateReport(result *BuildResult) string {
	if result == nil {
		return ""
	}

	report := BuildReport{
		RunID:       result.RunID,
		RunNumber:   result.RunNumber,
		Status:      result.Status,
		Conclusion:  result.Conclusion,
		Duration:    time.Now().Format(time.RFC3339),
		WorkflowURL: result.WorkflowURL,
		Artifact:    result.Artifact,
	}

	distDir := "dist"
	os.MkdirAll(distDir, 0755)

	jsonPath := filepath.Join(distDir, "build-report.json")
	data, _ := json.MarshalIndent(report, "", "  ")
	os.WriteFile(jsonPath, data, 0644)

	mdPath := filepath.Join(distDir, "build-report.md")
	mdContent := fmt.Sprintf(`# Build Report

| Field | Value |
|-------|-------|
| Run # | %d |
| Status | %s |
| Conclusion | %s |
| Artifact | %s |
| Workflow | [View](%s) |
`, result.RunNumber, result.Status, result.Conclusion, result.Artifact, result.WorkflowURL)
	os.WriteFile(mdPath, []byte(mdContent), 0644)

	fmt.Printf("Report: %s\n", jsonPath)
	return jsonPath
}
