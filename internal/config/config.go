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

type FlutterSettings struct {
	Enabled bool   `json:"enabled"`
	Channel string `json:"channel"`
}

type ReactNativeSettings struct {
	Enabled bool   `json:"enabled"`
	Entry   string `json:"entry_file"`
}

type RepoConfig struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

type GitHubConfig struct {
	Authenticated bool `json:"authenticated"`
}

type Config struct {
	ProjectName     string              `json:"project_name"`
	Repository      string              `json:"repository"`
	Repo            RepoConfig          `json:"repo"`
	GitHub          GitHubConfig        `json:"github"`
	IOS             iOSSettings         `json:"ios"`
	Signing         SigningSettings     `json:"signing"`
	MobAI           MobAISettings       `json:"mob_ai"`
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
		IOS: iOSSettings{
			MinimumVersion: "15.0",
			TargetVersion:  "17.0",
			Devices:        []string{"iPhone", "iPad"},
		},
		Signing: SigningSettings{},
		MobAI: MobAISettings{
			Enabled: false,
		},
		Flutter: FlutterSettings{
			Enabled: false,
			Channel: "stable",
		},
		ReactNative: ReactNativeSettings{
			Enabled: false,
			Entry:   "index.js",
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
	return errs
}
