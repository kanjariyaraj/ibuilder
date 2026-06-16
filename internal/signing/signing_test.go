package signing

import (
	"testing"
)

func TestValidateCertFile_Missing(t *testing.T) {
	err := ValidateCertFile("/nonexistent/cert.p12")
	if err == nil {
		t.Error("expected error for missing cert file")
	}
}

func TestValidateProfileFile_Missing(t *testing.T) {
	err := ValidateProfileFile("/nonexistent/profile.mobileprovision")
	if err == nil {
		t.Error("expected error for missing profile file")
	}
}

func TestValidateSigning_EmptyConfig(t *testing.T) {
	cfg := &SigningConfig{}
	result := ValidateSigning(nil, "", "", cfg)
	if result.Valid {
		t.Error("expected invalid result for empty config")
	}
}

func TestRunDoctor_EmptyConfig(t *testing.T) {
	cfg := &SigningConfig{}
	result := RunDoctor(nil, "", "", cfg)
	if result.Healthy {
		t.Error("expected unhealthy result for empty config")
	}
}

func TestGetStatus_EmptyConfig(t *testing.T) {
	status := GetStatus(nil, "", "", nil)
	if status.Configured {
		t.Error("expected not configured for nil config")
	}
}

func TestEncryptSecret(t *testing.T) {
	result, err := encryptSecret("dGVzdGtleQ==", "key123", "myvalue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.KeyID != "key123" {
		t.Errorf("expected key123, got %s", result.KeyID)
	}
	if result.EncryptedValue == "" {
		t.Errorf("expected non-empty encrypted value")
	}
}
