// Package gofraerror contains error extended which could be used as controller response next to error
package gofraerror

import "encoding/json"

const (
	ErrorTypePlainText ErrorType = iota
	ErrorTypeJSON
)

type ErrorType int

// CustomError represents a structured error with additional metadata
type Error struct {
	Message        string
	Type           ErrorType
	ResponseStatus int
}

func New(message string, errorType ErrorType, responseStatus int) *Error {
	return &Error{
		Message:        message,
		Type:           errorType,
		ResponseStatus: responseStatus,
	}
}

func NewText(message string, responseStatus int) *Error {
	return &Error{
		Message:        message,
		Type:           ErrorTypePlainText,
		ResponseStatus: responseStatus,
	}
}

func NewJSON(message string, responseStatus int) *Error {
	return &Error{
		Message:        message,
		Type:           ErrorTypeJSON,
		ResponseStatus: responseStatus,
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Type == ErrorTypePlainText {
		return e.Message
	}

	errorResponse := map[string]any{
		"error":  e.Message,
		"status": e.ResponseStatus,
	}

	// Not checking error as this should always work
	res, _ := json.Marshal(errorResponse)

	return string(res)
}
