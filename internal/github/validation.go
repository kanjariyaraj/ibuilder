package github

import "github.com/kanjariyaraj/Builder/internal/errors"

type ValidationResult struct {
	Authenticated  bool     `json:"authenticated"`
	RepoExists     bool     `json:"repo_exists"`
	ActionsEnabled bool     `json:"actions_enabled"`
	CanPush        bool     `json:"can_push"`
	CanAdmin       bool     `json:"can_admin"`
	Errors         []string `json:"errors,omitempty"`
}

func ValidateAll(token string, owner, name string) *ValidationResult {
	result := &ValidationResult{}

	if token == "" {
		result.Errors = append(result.Errors, "not authenticated")
		return result
	}

	if err := ValidateToken(token); err != nil {
		result.Errors = append(result.Errors, "authentication invalid: "+err.Error())
		return result
	}
	result.Authenticated = true

	client := NewClient(token)
	info, err := FetchRepoInfo(client, owner, name)
	if err != nil {
		if errors.IsKind(err, errors.KindNotFound) {
			result.Errors = append(result.Errors, "repository not found or no access")
		} else {
			result.Errors = append(result.Errors, "failed to validate repository: "+err.Error())
		}
		return result
	}

	result.RepoExists = true
	result.ActionsEnabled = info.ActionsEnabled
	result.CanPush = info.Permissions.Push
	result.CanAdmin = info.Permissions.Admin

	return result
}
