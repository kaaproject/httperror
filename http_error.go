// httperror provides a golang error-compatible type for escalating HTTP status codes alongside with the cause descriptions.
package httperror

import (
	"fmt"
	"net/http"
)

// HTTPError is a convenience error type for returning function processing results back to callers.
type HTTPError struct {
	statusCode  int
	description string
}

// Error function transforms HTTPError into a human-readable string.
func (p *HTTPError) Error() string {
	return p.description
}

// New constructs a new HTTPError.
func New(code int, format string, a ...interface{}) error {
	return &HTTPError{
		statusCode:  code,
		description: fmt.Sprintf(format, a...),
	}
}

// StatusCode is a convenience function for extracting HTTP Status Code from error types.
// It returns 200 for nil errors, embedded StatusCode for HTTPError instances, and 500 in every other case.
func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if err, ok := err.(*HTTPError); ok {
		return err.statusCode
	}

	return http.StatusInternalServerError
}

// ReasonPhrase is a convenience function for extracting HTTP Reason Phrase from error types.
func ReasonPhrase(err error) string {
	return http.StatusText(StatusCode(err))
}
