package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kanjariyaraj/Builder/internal/errors"
)

type RepoInfo struct {
	Owner          string `json:"owner"`
	Name           string `json:"name"`
	FullName       string `json:"full_name"`
	DefaultBranch  string `json:"default_branch"`
	Visibility     string `json:"visibility"`
	Permissions    Permissions `json:"permissions"`
	ActionsEnabled bool   `json:"actions_enabled"`
}

type Permissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type gitRemote struct {
	Owner string
	Name  string
}

func DetectGitRemote() (*gitRemote, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(errors.KindNotFound, "no git remote 'origin' found", err)
	}

	url := strings.TrimSpace(string(output))
	return parseRemoteURL(url)
}

func parseRemoteURL(url string) (*gitRemote, error) {
	if strings.HasPrefix(url, "https://") {
		parts := strings.Split(strings.TrimPrefix(url, "https://"), "/")
		if len(parts) >= 2 {
			owner := parts[len(parts)-2]
			name := strings.TrimSuffix(parts[len(parts)-1], ".git")
			return &gitRemote{Owner: owner, Name: name}, nil
		}
	}

	if strings.HasPrefix(url, "git@") {
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			owner := parts[len(parts)-2]
			name := strings.TrimSuffix(parts[len(parts)-1], ".git")
			return &gitRemote{Owner: owner, Name: name}, nil
		}
	}

	return nil, errors.New(errors.KindConfig, "unable to parse git remote URL")
}

type gitHubRepo struct {
	FullName      string                 `json:"full_name"`
	DefaultBranch string                 `json:"default_branch"`
	Visibility    string                 `json:"visibility"`
	Permissions   map[string]bool        `json:"permissions"`
}

func FetchRepoInfo(client *Client, owner, name string) (*RepoInfo, error) {
	data, err := client.Get(fmt.Sprintf("/repos/%s/%s", owner, name))
	if err != nil {
		return nil, err
	}

	var ghRepo gitHubRepo
	if err := json.Unmarshal(data, &ghRepo); err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to parse repo info", err)
	}

	info := &RepoInfo{
		Owner:         owner,
		Name:          name,
		FullName:      ghRepo.FullName,
		DefaultBranch: ghRepo.DefaultBranch,
		Visibility:    ghRepo.Visibility,
		Permissions: Permissions{
			Admin: ghRepo.Permissions["admin"],
			Push:  ghRepo.Permissions["push"],
			Pull:  ghRepo.Permissions["pull"],
		},
	}

	actionsEnabled, err := checkActionsEnabled(client, owner, name)
	if err == nil {
		info.ActionsEnabled = actionsEnabled
	}

	return info, nil
}

func checkActionsEnabled(client *Client, owner, name string) (bool, error) {
	_, err := client.Get(fmt.Sprintf("/repos/%s/%s/actions/workflows", owner, name))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func ValidateRepoAccess(client *Client, owner, name string) error {
	_, err := client.Get(fmt.Sprintf("/repos/%s/%s", owner, name))
	return err
}
