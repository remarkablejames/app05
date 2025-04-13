package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// Response represents a standardized API response envelope
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`    // Machine-readable error code
	Message string `json:"message"` // User-friendly error message
}

// MetaInfo contains metadata about the response
type MetaInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	PageInfo  *PageInfo `json:"pagination,omitempty"`
}

// PageInfo contains pagination information
type PageInfo struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalRows int `json:"total_rows"`
}

// APIVersion represents the current API version
const APIVersion = "v1"

// MaxRequestSize represents the maximum allowed request body size (1MB)
const MaxRequestSize = 1_048_576

// SendJSON sends a successful JSON response
func SendJSON(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
		Error:   nil,
		Meta: &MetaInfo{
			Timestamp: time.Now(),
			Version:   APIVersion,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If we fail to encode a success response, we have a serious problem
		panic(fmt.Sprintf("failed to encode success response: %v", err))
	}
}

// SendJSONWithPagination sends a successful JSON response with pagination info
func SendJSONWithPagination(w http.ResponseWriter, data interface{}, page, perPage, totalRows int) {
	response := Response{
		Success: true,
		Data:    data,
		Error:   nil,
		Meta: &MetaInfo{
			Timestamp: time.Now(),
			Version:   APIVersion,
			PageInfo: &PageInfo{
				Page:      page,
				PerPage:   perPage,
				TotalRows: totalRows,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(fmt.Sprintf("failed to encode paginated response: %v", err))
	}
}

// ParseJSON parses JSON from request body with size limits and strict parsing
func ParseJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		case errors.Is(err, http.ErrHandlerTimeout):
			return errors.New("request timeout")
		default:
			return err
		}
	}

	if dec.More() {
		return errors.New("body must contain a single JSON object")
	}

	return nil
}

// ParseUUID parses a UUID string

func ParseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID")
	}
	return id, nil
}
