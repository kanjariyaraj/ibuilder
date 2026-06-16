package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
)

type ProjectType string

const (
	ProjectNativeiOS     ProjectType = "native-ios"
	ProjectFlutter       ProjectType = "flutter"
	ProjectReactNative   ProjectType = "react-native"
	ProjectExpo          ProjectType = "expo"
	ProjectCapacitor     ProjectType = "capacitor"
	ProjectIonic         ProjectType = "ionic"
	ProjectCordova       ProjectType = "cordova"
	ProjectKMM           ProjectType = "kotlin-multiplatform"
	ProjectUnity         ProjectType = "unity"
	ProjectUnreal        ProjectType = "unreal"
	ProjectUnknown       ProjectType = "unknown"
)

type InitOptions struct {
	Force  bool
	DryRun bool
	Yes    bool
	JSON   bool
	Dir    string
}

type InitResult struct {
	ProjectType   ProjectType `json:"project_type"`
	ProjectName   string      `json:"project_name"`
	ConfigPath    string      `json:"config_path"`
	WorkflowPath  string      `json:"workflow_path"`
	TemplateUsed  string      `json:"template_used"`
	FilesCreated  []string    `json:"files_created"`
}

type Detector interface {
	Detect(dir string) (ProjectType, bool)
}

type SchemeDetector interface {
	DetectSchemes(dir string) ([]string, error)
	DetectWorkspaces(dir string) ([]string, error)
}

func Run(log *logger.Logger, opts *InitOptions) (*InitResult, error) {
	log.Info("initializing project...")

	pt := DetectProjectType(opts.Dir)
	if pt == ProjectUnknown {
		return nil, fmt.Errorf("unable to detect project type in %s", opts.Dir)
	}
	log.Info("detected project type: %s", string(pt))

	projectName := detectProjectName(opts.Dir, pt)
	log.Info("detected project name: %s", projectName)

	repoOwner, repoName := detectRepoInfo(opts.Dir)

	cfg := generateConfig(pt, projectName, repoOwner, repoName, opts.Dir)
	cfgPath := filepath.Join(opts.Dir, "builder.json")

	if opts.DryRun {
		showDryRun(pt, projectName, cfgPath, opts.Dir)
		return &InitResult{
			ProjectType:  pt,
			ProjectName:  projectName,
			ConfigPath:   cfgPath,
			TemplateUsed: string(pt),
		}, nil
	}

	if !opts.Force {
		if fileExists(cfgPath) {
			if !opts.Yes {
				fmt.Printf("Overwrite %s? [y/N]: ", cfgPath)
				var resp string
				fmt.Scanln(&resp)
				if resp != "y" && resp != "Y" {
					log.Info("skipped config generation")
					return nil, fmt.Errorf("cancelled by user")
				}
			}
		}
	}

	if err := config.Save(cfgPath, cfg); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}
	log.Info("created %s", cfgPath)

	workflowPath, err := generateWorkflow(pt, projectName, opts.Dir, opts.Force, opts.Yes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate workflow: %w", err)
	}
	log.Info("created %s", workflowPath)

	filesCreated := []string{cfgPath}
	if workflowPath != "" {
		filesCreated = append(filesCreated, workflowPath)
	}

	result := &InitResult{
		ProjectType:  pt,
		ProjectName:  projectName,
		ConfigPath:   cfgPath,
		WorkflowPath: workflowPath,
		TemplateUsed: string(pt),
		FilesCreated: filesCreated,
	}

	log.Info("project initialized successfully")
	return result, nil
}

func showDryRun(pt ProjectType, projectName, cfgPath, dir string) {
	fmt.Println("Dry Run - Files to create:")
	fmt.Printf("  - %s\n", cfgPath)
	wfDir := filepath.Join(dir, ".github", "workflows")
	fmt.Printf("  - %s/ios-build.yml\n", wfDir)
	fmt.Println()
	fmt.Printf("Project Type:  %s\n", pt)
	fmt.Printf("Project Name:  %s\n", projectName)
	fmt.Printf("Template:      %s\n", string(pt))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
