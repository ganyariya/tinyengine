package core

import "fmt"

// EngineError represents an error that occurred in the engine
type EngineError struct {
	Component string
	Operation string
	Cause     error
}

func (e *EngineError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("engine %s: %s failed: %v", e.Component, e.Operation, e.Cause)
	}
	return fmt.Sprintf("engine %s: %s failed", e.Component, e.Operation)
}

func (e *EngineError) Unwrap() error {
	return e.Cause
}

// NewEngineError creates a new EngineError
func NewEngineError(component, operation string, cause error) *EngineError {
	return &EngineError{
		Component: component,
		Operation: operation,
		Cause:     cause,
	}
}

// Predefined error types
var (
	ErrApplicationNotSet = NewEngineError("core", "application not set", nil)
)