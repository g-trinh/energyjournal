package error

import "fmt"

// InputValidationError represents a validation error for user input.
type InputValidationError struct {
	Field   string
	Message string
}

func (e *InputValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// NewInputValidationError creates a new InputValidationError.
func NewInputValidationError(field, message string) *InputValidationError {
	return &InputValidationError{Field: field, Message: message}
}

// NotFoundError represents a resource not found error.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{Resource: resource, ID: id}
}

// RateLimitError represents a rate limiting error.
type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "rate limit exceeded"
}

// NewRateLimitError creates a new RateLimitError.
func NewRateLimitError(message string) *RateLimitError {
	return &RateLimitError{Message: message}
}

// UnknownError represents an unexpected internal error.
type UnknownError struct {
	Err error
}

func (e *UnknownError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

func (e *UnknownError) Unwrap() error {
	return e.Err
}

// NewUnknownError wraps an error as an UnknownError.
func NewUnknownError(err error) *UnknownError {
	return &UnknownError{Err: err}
}
