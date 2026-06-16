package signing

import (
	"fmt"
	"os"

	"github.com/kanjariyaraj/Builder/internal/github"
)

type ValidationResult struct {
	Valid         bool     `json:"valid"`
	CertExists    bool     `json:"cert_exists"`
	ProfileExists bool     `json:"profile_exists"`
	SecretsSet    bool     `json:"secrets_set"`
	Errors        []string `json:"errors,omitempty"`
}

func ValidateSigning(client *github.Client, owner, repo string, cfg *SigningConfig) *ValidationResult {
	result := &ValidationResult{Valid: true}

	secrets := []string{
		cfg.CertSecret,
		"IOS_CERT_PASSWORD",
		cfg.ProfileSecret,
	}

	allSecretsExist := true
	for _, s := range secrets {
		if s == "" {
			continue
		}
		exists := checkSecretExists(client, owner, repo, s)
		if !exists {
			result.Errors = append(result.Errors, fmt.Sprintf("secret %s is not set", s))
			allSecretsExist = false
		}
	}

	result.SecretsSet = allSecretsExist
	if !allSecretsExist {
		result.Valid = false
	}

	return result
}

func checkSecretExists(client *github.Client, owner, repo, name string) bool {
	if client == nil {
		return false
	}
	_, err := client.Get(fmt.Sprintf("/repos/%s/%s/actions/secrets/%s", owner, repo, name))
	return err == nil
}

func ValidateCertFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("certificate file not found: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("certificate file is empty")
	}
	return nil
}

func ValidateProfileFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("provisioning profile not found: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("provisioning profile is empty")
	}
	return nil
}

type DoctorResult struct {
	Healthy       bool     `json:"healthy"`
	CertOK        bool     `json:"cert_ok"`
	ProfileOK     bool     `json:"profile_ok"`
	SecretsOK     bool     `json:"secrets_ok"`
	TeamIDSet     bool     `json:"team_id_set"`
	Issues        []string `json:"issues,omitempty"`
}

func RunDoctor(client *github.Client, owner, repo string, cfg *SigningConfig) *DoctorResult {
	result := &DoctorResult{Healthy: true}

	if cfg.TeamID == "" {
		result.Issues = append(result.Issues, "team_id is not configured in builder.json")
		result.Healthy = false
	} else {
		result.TeamIDSet = true
	}

	certSecret := cfg.CertSecret
	if certSecret == "" {
		certSecret = "IOS_CERT_BASE64"
	}

	if checkSecretExists(client, owner, repo, certSecret) {
		result.CertOK = true
	} else {
		result.Issues = append(result.Issues, fmt.Sprintf("secret %s not found", certSecret))
		result.Healthy = false
	}

	profileSecret := cfg.ProfileSecret
	if profileSecret == "" {
		profileSecret = "IOS_PROFILE_BASE64"
	}

	if checkSecretExists(client, owner, repo, profileSecret) {
		result.ProfileOK = true
	} else {
		result.Issues = append(result.Issues, fmt.Sprintf("secret %s not found", profileSecret))
		result.Healthy = false
	}

	result.SecretsOK = result.CertOK && result.ProfileOK

	return result
}
