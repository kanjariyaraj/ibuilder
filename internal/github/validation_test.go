package github

import (
	"testing"
)

func TestValidateAllNoToken(t *testing.T) {
	result := ValidateAll("", "owner", "repo")
	if result.Authenticated {
		t.Errorf("expected not authenticated")
	}
	if len(result.Errors) == 0 {
		t.Errorf("expected errors for no token")
	}
}
