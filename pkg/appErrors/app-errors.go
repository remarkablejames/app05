package appErrors

import (
	"fmt"
	"net/http"
	"runtime"
)

type Severity int

const (
	SeverityLow Severity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

type ErrorCode struct {
	Status   int
	Code     string
	Message  string
	Severity Severity
}

var (
	CodeInternal = ErrorCode{
		Status:   http.StatusInternalServerError,
		Code:     http.StatusText(http.StatusInternalServerError),
		Message:  "An internal error occurred",
		Severity: SeverityCritical,
	}

	CodeBadRequest = ErrorCode{
		Status:   http.StatusBadRequest,
		Code:     http.StatusText(http.StatusBadRequest),
		Message:  "Invalid request",
		Severity: SeverityLow,
	}

	CodeUnauthorized = ErrorCode{
		Status:   http.StatusUnauthorized,
		Code:     http.StatusText(http.StatusUnauthorized),
		Message:  "Authentication required",
		Severity: SeverityMedium,
	}

	CodeNotFound = ErrorCode{
		Status:   http.StatusNotFound,
		Code:     http.StatusText(http.StatusNotFound),
		Message:  "Resource not found",
		Severity: SeverityLow,
	}

	CodeForbidden = ErrorCode{
		Status:   http.StatusForbidden,
		Code:     http.StatusText(http.StatusForbidden),
		Message:  "Access denied",
		Severity: SeverityMedium,
	}
)

type AppError struct {
	Err     error
	Code    ErrorCode
	Message string
	Data    map[string]interface{}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code.Code, e.Message)
}

func New(code ErrorCode, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Data:    make(map[string]interface{}),
	}
}

func Wrap(err error, code ErrorCode) *AppError {
	if err == nil {
		return nil
	}

	appErr, ok := err.(*AppError)
	if ok {
		return appErr
	}

	return &AppError{
		Err:     err,
		Code:    code,
		Message: code.Message,
		Data:    make(map[string]interface{}),
	}
}

func getStack() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var stack string
	for {
		frame, more := frames.Next()
		stack += fmt.Sprintf("\n%s:%d", frame.Function, frame.Line)
		if !more {
			break
		}
	}
	return stack
}
