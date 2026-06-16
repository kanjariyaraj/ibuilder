package github

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kanjariyaraj/Builder/internal/errors"
)

type TokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func tokenDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(errors.KindInternal, "cannot determine home directory", err)
	}
	return filepath.Join(home, TokenDir), nil
}

func tokenPath() (string, error) {
	dir, err := tokenDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, TokenFile), nil
}

func LoadToken() (*TokenData, error) {
	path, err := tokenPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(errors.KindConfig, "failed to read token file", err)
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, errors.Wrap(errors.KindConfig, "corrupt token file", err)
	}

	if token.AccessToken == "" {
		return nil, nil
	}

	return &token, nil
}

func SaveToken(token *TokenData) error {
	dir, err := tokenDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.Wrap(errors.KindPermission, "failed to create token directory", err)
	}

	path, err := tokenPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return errors.Wrap(errors.KindInternal, "failed to marshal token", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return errors.Wrap(errors.KindPermission, "failed to write token file", err)
	}

	return nil
}

func DeleteToken() error {
	path, err := tokenPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.Wrap(errors.KindPermission, "failed to delete token file", err)
	}

	return nil
}

func TokenExists() bool {
	token, err := LoadToken()
	return err == nil && token != nil
}
