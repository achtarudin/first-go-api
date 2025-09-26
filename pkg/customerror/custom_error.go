package customerror

import (
	"fmt"
	"net/http"
)

// Define some common error types
const (
	CodeValidationFailed = "VALIDATION_FAILED"
	CodeInvalidInput     = "INVALID_INPUT"
	CodeUnauthorized     = "UNAUTHORIZED"
	CodeHashFailed       = "HASH_PASSWORD_FAILED"
	CodeDBQuery          = "DB_QUERY_FAILED"
	CodeNotFound         = "NOT_FOUND"
	CodeConflict         = "CONFLICT"
	CodeInternal         = "INTERNAL_ERROR"
)

// Struct for application-specific errors
type CustomError struct {
	code    string
	message string
	fields  map[string]any
	cause   error
}

// Implement the error interface
func (e *CustomError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

// Unwrap allows errors.Is and errors.As to work with CustomeError
func (e *CustomError) Unwrap() error { return e.cause }

// Accessor methods
func (e *CustomError) Message() string {
	return e.message
}

func (e *CustomError) Code() string {
	return e.code
}

func (e *CustomError) Fields() map[string]any {
	return cloneMap(e.fields)
}

func (e *CustomError) Cause() error {
	if e.cause == nil {
		return nil
	}
	return e.cause
}

func (e *CustomError) StatusCode() int {
	switch e.code {
	case CodeInvalidInput:
		return http.StatusBadRequest
	case CodeValidationFailed:
		return http.StatusUnprocessableEntity
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeHashFailed, CodeDBQuery, CodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Constructor for CustomError
func New(code string, message string, fields map[string]any, cause error) *CustomError {
	return &CustomError{
		code:    code,
		message: message,
		fields:  cloneMap(fields),
		cause:   cause,
	}
}

func NewFieldToAny(code string, message string, fields map[string]string, cause error) *CustomError {

	return New(code, message, cloneMapString(fields), cause)
}

// Helper to clone a map to avoid external modifications
func cloneMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func cloneMapString(src map[string]string) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
