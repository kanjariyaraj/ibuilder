package signing

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
)

type SigningConfig struct {
	TeamID        string `json:"team_id"`
	BundleID      string `json:"bundle_id"`
	CertSecret    string `json:"certificate_secret"`
	ProfileSecret string `json:"profile_secret"`
}

type CertInfo struct {
	Path       string
	Password   string
	TeamID     string
	BundleID   string
	Encoded    string
	Hash       string
}

type ProfileInfo struct {
	Path         string
	TeamID       string
	BundleID     string
	Expiration   string
	Encoded      string
	Hash         string
}

type SetupOptions struct {
	CertPath     string
	CertPassword string
	ProfilePath  string
	TeamID       string
	BundleID     string
	Force        bool
}

type SetupResult struct {
	CertSecret    string `json:"cert_secret"`
	ProfileSecret string `json:"profile_secret"`
	TeamID        string `json:"team_id"`
	BundleID      string `json:"bundle_id"`
}

func RunSetup(client *github.Client, opts *SetupOptions, cfg *config.Config) (*SetupResult, error) {
	certData, err := os.ReadFile(opts.CertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	profileData, err := os.ReadFile(opts.ProfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read provisioning profile: %w", err)
	}

	certEncoded := base64.StdEncoding.EncodeToString(certData)
	profileEncoded := base64.StdEncoding.EncodeToString(profileData)

	certSecret := "IOS_CERT_BASE64"
	certPassSecret := "IOS_CERT_PASSWORD"
	profileSecret := "IOS_PROFILE_BASE64"

	if err := uploadSecret(client, cfg.Repo.Owner, cfg.Repo.Name, certSecret, certEncoded); err != nil {
		return nil, fmt.Errorf("failed to upload certificate secret: %w", err)
	}

	if err := uploadSecret(client, cfg.Repo.Owner, cfg.Repo.Name, certPassSecret, opts.CertPassword); err != nil {
		return nil, fmt.Errorf("failed to upload certificate password secret: %w", err)
	}

	if err := uploadSecret(client, cfg.Repo.Owner, cfg.Repo.Name, profileSecret, profileEncoded); err != nil {
		return nil, fmt.Errorf("failed to upload profile secret: %w", err)
	}

	if opts.TeamID != "" {
		if err := uploadSecret(client, cfg.Repo.Owner, cfg.Repo.Name, "IOS_TEAM_ID", opts.TeamID); err != nil {
			return nil, fmt.Errorf("failed to upload team ID secret: %w", err)
		}
	}

	if opts.BundleID != "" {
		if err := uploadSecret(client, cfg.Repo.Owner, cfg.Repo.Name, "IOS_BUNDLE_ID", opts.BundleID); err != nil {
			return nil, fmt.Errorf("failed to upload bundle ID secret: %w", err)
		}
	}

	cfg.Signing.TeamID = opts.TeamID

	result := &SetupResult{
		CertSecret:    certSecret,
		ProfileSecret: profileSecret,
		TeamID:        opts.TeamID,
		BundleID:      opts.BundleID,
	}

	return result, nil
}

func uploadSecret(client *github.Client, owner, repo, name, value string) error {
	path := fmt.Sprintf("/repos/%s/%s/actions/secrets/%s", owner, repo, name)

	publicKey, err := getPublicKey(client, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	encrypted, err := encryptSecret(publicKey.Key, publicKey.KeyID, value)
	if err != nil {
		return err
	}

	body, _ := json.Marshal(encrypted)
	_, err = client.Put(path, string(body))
	return err
}

type publicKeyResponse struct {
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
}

type encryptedSecret struct {
	EncryptedValue string `json:"encrypted_value"`
	KeyID          string `json:"key_id"`
}

func getPublicKey(client *github.Client, owner, repo string) (*publicKeyResponse, error) {
	data, err := client.Get(fmt.Sprintf("/repos/%s/%s/actions/secrets/public-key", owner, repo))
	if err != nil {
		return nil, err
	}

	var pk publicKeyResponse
	if err := json.Unmarshal(data, &pk); err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return &pk, nil
}

func encryptSecret(publicKey, keyID, value string) (*encryptedSecret, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	hash := sha256.Sum256(append(keyBytes, []byte(value)...))
	encrypted := base64.StdEncoding.EncodeToString(hash[:])

	return &encryptedSecret{
		EncryptedValue: encrypted,
		KeyID:          keyID,
	}, nil
}
