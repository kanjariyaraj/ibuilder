package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWorkflowStructs(t *testing.T) {
	w := Workflow{ID: 1, Name: "test", Path: ".github/workflows/test.yml"}
	if w.ID != 1 {
		t.Errorf("expected ID 1")
	}
	if w.Name != "test" {
		t.Errorf("expected Name 'test'")
	}
}

func TestWorkflowRunStruct(t *testing.T) {
	r := WorkflowRun{ID: 123, Status: "completed", Conclusion: "success"}
	if r.ID != 123 {
		t.Errorf("expected ID 123")
	}
	if r.Status != "completed" {
		t.Errorf("expected completed")
	}
}

func TestArtifactStruct(t *testing.T) {
	a := Artifact{ID: 1, Name: "build.ipa", Size: 1024}
	if a.Name != "build.ipa" {
		t.Errorf("expected 'build.ipa'")
	}
}

func TestListWorkflowsWithMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		resp := map[string][]Workflow{
			"workflows": {
				{ID: 1, Name: "CI", Path: ".github/workflows/ci.yml"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	origAPI := APIRoot
	APIRoot = ts.URL
	defer func() { APIRoot = origAPI }()

	client := NewClient("test-token")
	workflows, err := ListWorkflows(client, "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(workflows))
	}
}

func TestGetWorkflowRunWithMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := WorkflowRun{ID: 42, Status: "in_progress", Conclusion: ""}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	origAPI := APIRoot
	APIRoot = ts.URL
	defer func() { APIRoot = origAPI }()

	client := NewClient("test-token")
	run, err := GetWorkflowRun(client, "owner", "repo", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if run.ID != 42 {
		t.Errorf("expected ID 42, got %d", run.ID)
	}
	if run.Status != "in_progress" {
		t.Errorf("expected 'in_progress', got '%s'", run.Status)
	}
}

func TestListWorkflowRunsWithMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := workflowRunsResponse{
			TotalCount: 2,
			WorkflowRuns: []WorkflowRun{
				{ID: 1, Status: "completed", Conclusion: "success"},
				{ID: 2, Status: "in_progress", Conclusion: ""},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	origAPI := APIRoot
	APIRoot = ts.URL
	defer func() { APIRoot = origAPI }()

	client := NewClient("test-token")
	runs, err := ListWorkflowRuns(client, "owner", "repo", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 2 {
		t.Errorf("expected 2 runs, got %d", len(runs))
	}
}

func TestDispatchWorkflowWithMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	origAPI := APIRoot
	APIRoot = ts.URL
	defer func() { APIRoot = origAPI }()

	client := NewClient("test-token")
	err := DispatchWorkflow(client, "owner", "repo", "ci.yml", "main", map[string]string{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListArtifactsWithMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := artifactsResponse{
			Artifacts: []Artifact{
				{ID: 1, Name: "build.ipa", Size: 2048},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	origAPI := APIRoot
	APIRoot = ts.URL
	defer func() { APIRoot = origAPI }()

	client := NewClient("test-token")
	artifacts, err := ListArtifacts(client, "owner", "repo", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(artifacts) != 1 {
		t.Errorf("expected 1 artifact, got %d", len(artifacts))
	}
}
