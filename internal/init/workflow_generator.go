package init

import (
	"fmt"
	"os"
	"path/filepath"
)

func generateWorkflow(pt ProjectType, projectName, dir string, force, yes bool) (string, error) {
	wfDir := filepath.Join(dir, ".github", "workflows")
	if err := os.MkdirAll(wfDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create workflows directory: %w", err)
	}

	wfPath := filepath.Join(wfDir, "ios-build.yml")

	if fileExists(wfPath) && !force {
		if !yes {
			fmt.Printf("Overwrite %s? [y/N]: ", wfPath)
			var resp string
			fmt.Scanln(&resp)
			if resp != "y" && resp != "Y" {
				return "", nil
			}
		}
	}

	content := getTemplate(pt, projectName)
	if err := os.WriteFile(wfPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write workflow: %w", err)
	}

	return wfPath, nil
}
