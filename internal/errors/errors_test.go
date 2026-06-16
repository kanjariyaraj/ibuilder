package errors

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(KindValidation, "test error")
	if err.Kind != KindValidation {
		t.Errorf("expected KindValidation, got %s", err.Kind)
	}
	if err.Message != "test error" {
		t.Errorf("expected 'test error', got '%s'", err.Message)
	}
}

func TestNewError(t *testing.T) {
	err := New(KindValidation, "test message")
	expected := "validation: test message"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestNewWrappedError(t *testing.T) {
	inner := fmt.Errorf("inner error")
	err := Wrap(KindInternal, "wrapped", inner)
	expected := "internal: wrapped: inner error"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestUnwrap(t *testing.T) {
	inner := fmt.Errorf("inner error")
	err := Wrap(KindInternal, "wrapped", inner)
	unwrapped := err.Unwrap()
	if unwrapped != inner {
		t.Errorf("expected inner error, got %v", unwrapped)
	}
}

func TestIsKind(t *testing.T) {
	err := New(KindConfig, "config error")
	if !IsKind(err, KindConfig) {
		t.Errorf("expected true for KindConfig")
	}
	if IsKind(err, KindValidation) {
		t.Errorf("expected false for KindValidation")
	}
}

func TestIsKindNil(t *testing.T) {
	if IsKind(nil, KindConfig) {
		t.Errorf("expected false for nil error")
	}
}
