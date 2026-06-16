package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kanjariyaraj/Builder/internal/errors"
)

type iOSSettings struct {
	MinimumVersion string   `json:"minimum_version"`
	TargetVersion  string   `json:"target_version"`
	Devices        []string `json:"devices"`
}

type SigningSettings struct {
	TeamID        string `json:"team_id"`
	Provisioning  string `json:"provisioning_profile"`
	Certificate   string `json:"certificate"`
}

type MobAISettings struct {
	Enabled bool   `json:"enabled"`
	APIKey  string `json:"api_key"`
}

type MobaiSettings struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Device            string `json:"device"`
	AutoReconnect     bool   `json:"auto_reconnect"`
	ConnectionTimeout int    `json:"connection_timeout"`
}

type AISettings struct {
	Enabled             bool    `json:"enabled"`
	AutoFix             bool    `json:"auto_fix"`
	ReportFormat        string  `json:"report_format"`
	ConfidenceThreshold float64 `json:"confidence_threshold"`
}

type FlutterSettings struct {
	Enabled     bool   `json:"enabled"`
	Channel     string `json:"channel"`
	Watch       bool   `json:"watch"`
	HotReload   bool   `json:"hot_reload"`
	DebounceMs  int    `json:"debounce_ms"`
	AutoAttach  bool   `json:"auto_attach"`
	AutoInstall bool   `json:"auto_install"`
}

type ReactNativeSettings struct {
	Enabled        bool   `json:"enabled"`
	Entry          string `json:"entry_file"`
	MetroPort      int    `json:"metro_port"`
	AutoStartMetro bool   `json:"auto_start_metro"`
	AutoAttach     bool   `json:"auto_attach"`
	AutoInstall    bool   `json:"auto_install"`
	FastRefresh    bool   `json:"fast_refresh"`
}

type ReleaseSettings struct {
	Mode                 string `json:"mode"`
	AutoNotes            bool   `json:"auto_notes"`
	DefaultGroup         string `json:"default_group"`
	ValidateBeforeUpload bool   `json:"validate_before_upload"`
	CreateGitHubRelease  bool   `json:"create_github_release"`
	Notifications        bool   `json:"notifications"`
}

type RepoConfig struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

type GitHubConfig struct {
	Authenticated bool `json:"authenticated"`
}

type BuildConfig struct {
	WorkflowID    string `json:"workflow_id"`
	Branch        string `json:"branch"`
	Scheme        string `json:"scheme"`
	Configuration string `json:"configuration"`
	BuildMode     string `json:"build_mode"`
	ProjectType   string `json:"project_type"`
}

type Config struct {
	ProjectName     string              `json:"project_name"`
	Repository      string              `json:"repository"`
	Repo            RepoConfig          `json:"repo"`
	GitHub          GitHubConfig        `json:"github"`
	Build           BuildConfig         `json:"build"`
	IOS             iOSSettings         `json:"ios"`
	Signing         SigningSettings     `json:"signing"`
	MobAI           MobAISettings       `json:"mob_ai"`
	Mobai           MobaiSettings       `json:"mobai"`
	AI              AISettings          `json:"ai"`
	Release         ReleaseSettings     `json:"release"`
	Flutter         FlutterSettings     `json:"flutter"`
	ReactNative     ReactNativeSettings `json:"react_native"`
}

func Default() *Config {
	return &Config{
		ProjectName: "Builder",
		Repository:  "",
		Repo: RepoConfig{
			Owner:  "",
			Name:   "",
			Branch: "main",
		},
		GitHub: GitHubConfig{
			Authenticated: false,
		},
		Build: BuildConfig{
			WorkflowID:    "",
			Branch:        "main",
			Scheme:        "",
			Configuration: "Release",
			BuildMode:     "release",
			ProjectType:   "xcode",
		},
		IOS: iOSSettings{
			MinimumVersion: "15.0",
			TargetVersion:  "17.0",
			Devices:        []string{"iPhone", "iPad"},
		},
		Signing: SigningSettings{},
		MobAI: MobAISettings{
			Enabled: false,
		},
		Mobai: MobaiSettings{
			Host:              "",
			Port:              0,
			Device:            "",
			AutoReconnect:     true,
			ConnectionTimeout: 30,
		},
		AI: AISettings{
			Enabled:             true,
			AutoFix:             false,
			ReportFormat:        "markdown",
			ConfidenceThreshold: 0.8,
		},
		Release: ReleaseSettings{
			Mode:                 "beta",
			AutoNotes:            true,
			DefaultGroup:         "",
			ValidateBeforeUpload: true,
			CreateGitHubRelease:  true,
			Notifications:        false,
		},
		Flutter: FlutterSettings{
			Enabled:     false,
			Channel:     "stable",
			Watch:       true,
			HotReload:   true,
			DebounceMs:  500,
			AutoAttach:  true,
			AutoInstall: true,
		},
		ReactNative: ReactNativeSettings{
			Enabled:        false,
			Entry:          "index.js",
			MetroPort:      8081,
			AutoStartMetro: true,
			AutoAttach:     true,
			AutoInstall:    true,
			FastRefresh:    true,
		},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, errors.Wrap(errors.KindConfig, "failed to read config file", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Wrap(errors.KindConfig, "failed to parse config file", err)
	}
	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrap(errors.KindPermission, "failed to create config directory", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return errors.Wrap(errors.KindInternal, "failed to marshal config", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return errors.Wrap(errors.KindPermission, "failed to write config file", err)
	}
	return nil
}

func Validate(cfg *Config) []error {
	var errs []error
	if cfg.ProjectName == "" {
		errs = append(errs, errors.New(errors.KindValidation, "project_name is required"))
	}
	if cfg.IOS.MinimumVersion == "" {
		errs = append(errs, errors.New(errors.KindValidation, "ios.minimum_version is required"))
	}
	if cfg.IOS.TargetVersion == "" {
		errs = append(errs, errors.New(errors.KindValidation, "ios.target_version is required"))
	}
	if len(cfg.IOS.Devices) == 0 {
		errs = append(errs, errors.New(errors.KindValidation, "ios.devices must have at least one device"))
	}
	if cfg.Repo.Branch == "" {
		errs = append(errs, errors.New(errors.KindValidation, "repo.branch is required"))
	}
	if cfg.Build.Branch == "" {
		errs = append(errs, errors.New(errors.KindValidation, "build.branch is required"))
	}
	if cfg.Build.BuildMode != "debug" && cfg.Build.BuildMode != "release" {
		errs = append(errs, errors.New(errors.KindValidation, "build.build_mode must be 'debug' or 'release'"))
	}
	return errs
}
