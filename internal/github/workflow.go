package github

import "github.com/kanjariyaraj/Builder/internal/errors"

type Workflow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func ListWorkflows(client *Client, owner, name string) ([]Workflow, error) {
	_, err := client.Get("/repos/" + owner + "/" + name + "/actions/workflows")
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "failed to list workflows", err)
	}
	return nil, nil
}
