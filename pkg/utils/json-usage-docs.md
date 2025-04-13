# HTTP Response Utilities Documentation

This document provides comprehensive documentation for the HTTP response utilities package, which offers standardized JSON response handling, error management, and request parsing for Go web applications.

## Table of Contents
- [Response Structures](#response-structures)
- [Core Functions](#core-functions)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Response Structures

### Standard Response Envelope
All API responses follow a consistent structure:

```json
{
    "success": true|false,
    "data": { ... },      // (optional) Present for successful responses
    "error": { ... },     // (optional) Present for error responses
    "meta": { ... }       // Always present
}
```

### Meta Information
Every response includes metadata:
```json
{
    "meta": {
        "timestamp": "2025-02-08T12:34:56Z",
        "request_id": "req_abc123",
        "version": "v1",
        "processed_in": "1.234ms",
        "pagination": {    // (optional) Present for paginated responses
            "page": 1,
            "per_page": 10,
            "total_rows": 100
        }
    }
}
```

## Core Functions

### SendJSON
Sends a successful JSON response with the standard envelope.

```go
func someHandler(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()
    reqID := generateRequestID()
    
    data := struct {
        Name string `json:"name"`
    }{
        Name: "Example",
    }
    
    utils.SendJSON(w, data, reqID, startTime)
}
```

### SendJSONWithPagination
Sends a paginated JSON response.

```go
func listHandler(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()
    reqID := generateRequestID()
    
    items := []Item{...}
    
    utils.SendJSONWithPagination(
        w,
        items,
        page,      // current page number
        perPage,   // items per page
        totalRows, // total number of items
        reqID,
        startTime,
    )
}
```

### ParseJSON
Parses and validates JSON request bodies with size limits.

```go
func createHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    if err := utils.ParseJSON(w, r, &input); err != nil {
        // Handle parsing error
        return
    }
    
    // Process input...
}
```

## Error Handling

### SendError
Base function for sending error responses.

```go
utils.SendError(
    w,
    http.StatusBadRequest,
    "BAD_REQUEST",
    "Invalid input provided",
    "Field 'email' must be a valid email address",
    reqID,
    logger,
)
```

### Convenience Error Functions

#### SendInternalError
Use for internal server errors (500):
```go
if err != nil {
    utils.SendInternalError(w, r, err, reqID, logger)
    return
}
```

#### SendBadRequest
Use for client errors (400):
```go
if err := validate(input); err != nil {
    utils.SendBadRequest(w, r, err, reqID, logger)
    return
}
```

#### SendNotFound
Use when resources aren't found (404):
```go
if user == nil {
    utils.SendNotFound(w, r, reqID, logger)
    return
}
```

#### SendUnauthorized
Use for authentication errors (401):
```go
if !isAuthenticated(r) {
    utils.SendUnauthorized(w, r, reqID, logger)
    return
}
```

## Examples

### Complete Handler Example
```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()
    reqID := generateRequestID()
    logger := getLogger()

    // 1. Parse input
    var input struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    if err := utils.ParseJSON(w, r, &input); err != nil {
        utils.SendBadRequest(w, r, err, reqID, logger)
        return
    }
    
    // 2. Validate input
    if err := validateUser(input); err != nil {
        utils.SendBadRequest(w, r, err, reqID, logger)
        return
    }
    
    // 3. Process request
    user, err := createUser(input)
    if err != nil {
        utils.SendInternalError(w, r, err, reqID, logger)
        return
    }
    
    // 4. Send response
    utils.SendJSON(w, user, reqID, startTime)
}
```

### Paginated List Example
```go
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()
    reqID := generateRequestID()
    logger := getLogger()

    // Get pagination parameters
    page := getPaginationPage(r)
    perPage := getPaginationPerPage(r)
    
    // Fetch users
    users, total, err := getUsersList(page, perPage)
    if err != nil {
        utils.SendInternalError(w, r, err, reqID, logger)
        return
    }
    
    // Send paginated response
    utils.SendJSONWithPagination(w, users, page, perPage, total, reqID, startTime)
}
```

## Constants and Configuration

### Error Codes
```go
const (
    ErrCodeInternal     = "INTERNAL_ERROR"
    ErrCodeBadRequest   = "BAD_REQUEST"
    ErrCodeNotFound     = "NOT_FOUND"
    ErrCodeUnauthorized = "UNAUTHORIZED"
)
```

### Request Limits
```go
const MaxRequestSize = 1_048_576 // 1MB maximum request body size
```

### API Version
```go
const APIVersion = "v1"
```

## Best Practices

1. Always include request IDs in your handlers for tracing
2. Use appropriate error types for different scenarios
3. Include relevant logging with each error response
4. Set proper content-type headers (done automatically by these utilities)
5. Validate input early in your handlers
6. Use pagination for list endpoints
7. Include processing time for performance monitoring

## Notes

- All responses automatically include timestamps and request IDs
- JSON parsing automatically rejects unknown fields
- Request bodies are limited to 1MB by default
- All error responses include structured error information
- Pagination information is included when using `SendJSONWithPagination`