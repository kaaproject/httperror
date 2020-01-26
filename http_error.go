// httperror provides a golang error-compatible type for escalating HTTP status codes alongside with the cause descriptions.
// Copyright 2019 KaaIoT Technologies, LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httperror

import (
	"encoding/json"
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

// Write to the response writer a status code and a JSON-encoded message based on the provided error.
// The payload will look like:
// {
// 		"message": "error writing to DB"
// }
// In case the status code is 500 (Internal Server Error), the message will always be "Internal Server Error".
// Write does not end the request; the caller should ensure no further writes are done to w.
func Write(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	code := StatusCode(err)
	w.WriteHeader(code)

	var message string
	if code == http.StatusInternalServerError {
		message = http.StatusText(code)
	} else {
		message = err.Error()
	}

	_ = json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: message})
}
