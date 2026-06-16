package signing

import (
	"os"
	"testing"
)

func TestValidateCertFile_Valid(t *testing.T) {
	tmp, _ := os.CreateTemp("", "cert*.p12")
	defer os.Remove(tmp.Name())
	tmp.Write([]byte("test-cert-data"))
	tmp.Close()

	err := ValidateCertFile(tmp.Name())
	if err != nil {
		t.Errorf("expected no error for valid cert file, got %v", err)
	}
}

func TestValidateProfileFile_Valid(t *testing.T) {
	tmp, _ := os.CreateTemp("", "profile*.mobileprovision")
	defer os.Remove(tmp.Name())
	tmp.Write([]byte("test-profile-data"))
	tmp.Close()

	err := ValidateProfileFile(tmp.Name())
	if err != nil {
		t.Errorf("expected no error for valid profile file, got %v", err)
	}
}

func TestValidateCertFile_Empty(t *testing.T) {
	tmp, _ := os.CreateTemp("", "empty*.p12")
	defer os.Remove(tmp.Name())

	err := ValidateCertFile(tmp.Name())
	if err == nil {
		t.Error("expected error for empty cert file")
	}
}

func TestValidateProfileFile_Empty(t *testing.T) {
	tmp, _ := os.CreateTemp("", "empty*.mobileprovision")
	defer os.Remove(tmp.Name())

	err := ValidateProfileFile(tmp.Name())
	if err == nil {
		t.Error("expected error for empty profile file")
	}
}

func TestRunDoctor_TeamIDSet(t *testing.T) {
	cfg := &SigningConfig{TeamID: "TEAM123"}
	result := RunDoctor(nil, "", "", cfg)
	if !result.TeamIDSet {
		t.Error("expected TeamIDSet to be true")
	}
}

func TestRunDoctor_NoTeamID(t *testing.T) {
	cfg := &SigningConfig{}
	result := RunDoctor(nil, "", "", cfg)
	if result.TeamIDSet {
		t.Error("expected TeamIDSet to be false")
	}
	if result.Healthy {
		t.Error("expected unhealthy without team ID")
	}
}

func TestGetStatus_Configured(t *testing.T) {
	status := GetStatus(nil, "owner", "repo", nil)
	if status == nil {
		t.Error("expected non-nil status")
	}
}
