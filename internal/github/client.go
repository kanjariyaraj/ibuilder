package github

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/kanjariyaraj/Builder/internal/errors"
)

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "BuilderCLI/1.0")
	return c.httpClient.Do(req)
}

func (c *Client) AccessToken() string {
	return c.token
}

func (c *Client) Get(path string) ([]byte, error) {
	url := APIRoot + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to create request", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, fmt.Sprintf("API request failed: %s", path), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil, errors.New(errors.KindPermission, "authentication failed or token lacks permissions")
	}
	if resp.StatusCode == 404 {
		return nil, errors.New(errors.KindNotFound, "resource not found")
	}
	if resp.StatusCode == 429 {
		return nil, errors.New(errors.KindNetwork, "rate limit exceeded, please wait")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Wrap(errors.KindNetwork, fmt.Sprintf("API error (status %d)", resp.StatusCode), fmt.Errorf("%s", string(body)))
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) Put(path, body string) ([]byte, error) {
	url := APIRoot + path
	req, err := c.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, fmt.Sprintf("API request failed: %s", path), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil, errors.New(errors.KindPermission, "authentication failed or token lacks permissions")
	}
	if resp.StatusCode == 404 {
		return nil, errors.New(errors.KindNotFound, "resource not found")
	}
	if resp.StatusCode == 409 {
		return nil, errors.New(errors.KindNetwork, "conflict: resource already exists")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.Wrap(errors.KindNetwork, fmt.Sprintf("API error (status %d)", resp.StatusCode), fmt.Errorf("%s", string(bodyBytes)))
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) NewRequest(method, url, body string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func ValidateToken(token string) error {
	client := NewClient(token)
	_, err := client.Get("/user")
	return err
}
