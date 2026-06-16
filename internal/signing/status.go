package signing

import (
	"fmt"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
)

type SigningStatus struct {
	CertSecret      string `json:"cert_secret"`
	ProfileSecret   string `json:"profile_secret"`
	TeamID          string `json:"team_id"`
	SecretsUploaded bool   `json:"secrets_uploaded"`
	Configured      bool   `json:"configured"`
}

func GetStatus(client *github.Client, owner, repo string, cfg *config.Config) *SigningStatus {
	status := &SigningStatus{}

	if cfg != nil {
		status.TeamID = cfg.Signing.TeamID
		status.Configured = cfg.Signing.TeamID != ""
	}

	certSecret := "IOS_CERT_BASE64"
	profileSecret := "IOS_PROFILE_BASE64"

	if cfg != nil {
		if cfg.Signing.Certificate != "" {
			certSecret = cfg.Signing.Certificate
		}
		if cfg.Signing.Provisioning != "" {
			profileSecret = cfg.Signing.Provisioning
		}
	}

	status.CertSecret = certSecret
	status.ProfileSecret = profileSecret

	if client != nil && owner != "" && repo != "" {
		_, err := client.Get(fmt.Sprintf("/repos/%s/%s/actions/secrets/%s", owner, repo, certSecret))
		status.SecretsUploaded = err == nil
	}

	return status
}

func RemoveSigning(client *github.Client, owner, repo string) error {
	secrets := []string{"IOS_CERT_BASE64", "IOS_CERT_PASSWORD", "IOS_PROFILE_BASE64", "IOS_TEAM_ID", "IOS_BUNDLE_ID"}
	for _, s := range secrets {
		if err := deleteSecret(client, owner, repo, s); err != nil {
			return fmt.Errorf("failed to delete secret %s: %w", s, err)
		}
	}
	return nil
}

func deleteSecret(client *github.Client, owner, repo, name string) error {
	if client == nil {
		return nil
	}
	_, err := client.Get(fmt.Sprintf("/repos/%s/%s/actions/secrets/%s", owner, repo, name))
	if err != nil {
		return nil
	}
	return nil
}
