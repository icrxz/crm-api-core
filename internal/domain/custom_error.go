package domain

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	messagePrefix string
	message       string
	metadata      map[string]any
	statusCode    int
}

func (e CustomError) Error() string {
	return fmt.Sprintf("%s %s", e.messagePrefix, e.message)
}

func (e CustomError) Metadata() map[string]any {
	return e.metadata
}

func (e CustomError) StatusCode() int {
	return e.statusCode
}

func NewValidationError(message string, metadata map[string]any) error {
	return &CustomError{
		messagePrefix: "Validation error - Message:",
		message:       message,
		statusCode:    http.StatusBadRequest,
		metadata:      metadata,
	}
}

func NewParserError(message string, metadata map[string]any) error {
	return &CustomError{
		messagePrefix: "Parser error - Message:",
		message:       message,
		statusCode:    http.StatusBadRequest,
		metadata:      metadata,
	}
}

func NewNotFoundError(message string, metadata map[string]any) error {
	return &CustomError{
		messagePrefix: "NotFound error - Message:",
		message:       message,
		statusCode:    http.StatusNotFound,
		metadata:      metadata,
	}
}

func NewConflictError(message string, metadata map[string]any) error {
	return &CustomError{
		messagePrefix: "Conflict error - Message:",
		message:       message,
		statusCode:    http.StatusConflict,
		metadata:      metadata,
	}
}

func NewUnauthorizedError(message string) error {
	return &CustomError{
		messagePrefix: "Unauthorized error - Message:",
		message:       message,
		statusCode:    http.StatusUnauthorized,
	}
}

func (e CustomError) IsNotFound() bool {
	return e.statusCode == http.StatusNotFound
}
