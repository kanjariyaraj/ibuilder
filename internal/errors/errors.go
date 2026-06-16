package errors

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindConfig     Kind = "config"
	KindValidation Kind = "validation"
	KindNotFound   Kind = "not_found"
	KindPermission Kind = "permission"
	KindNetwork    Kind = "network"
	KindInternal   Kind = "internal"
)

type BuilderError struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *BuilderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func (e *BuilderError) Unwrap() error {
	return e.Err
}

func New(kind Kind, message string) *BuilderError {
	return &BuilderError{Kind: kind, Message: message}
}

func Wrap(kind Kind, message string, err error) *BuilderError {
	return &BuilderError{Kind: kind, Message: message, Err: err}
}

func IsKind(err error, kind Kind) bool {
	var bErr *BuilderError
	if As(err, &bErr) {
		return bErr.Kind == kind
	}
	return false
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
