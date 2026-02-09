package main

import (
	"fmt"
	"net/http"
)

// ErrorResponse represents the API error format specified in requirements
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// logError logs the error with details
func (b *backend) logError(r *http.Request, err error) {
	b.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

// errorResponse sends a JSON error response with the specified code and message
func (b *backend) errorResponse(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message

	err := b.writeJson(w, status, resp, nil)
	if err != nil {
		b.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse sends a 500 Internal Server Error response
func (b *backend) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	b.logError(r, err)
	message := fmt.Sprintf("the server encountered a problem and could not process your request")
	b.errorResponse(w, r, http.StatusInternalServerError, "SERVER_ERROR", message)
}

// notFoundResponse sends a 404 Not Found response
func (b *backend) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the requested resource could not be found: %s, %s", r.URL.Path, r.URL.RawQuery)
	b.errorResponse(w, r, http.StatusNotFound, "NOT_FOUND", message)
}

// badRequestResponse sends a 400 Bad Request response
func (b *backend) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	b.errorResponse(w, r, http.StatusBadRequest, "BAD_REQUEST", err.Error())
}

// validationErrorResponse sends a 400 Bad Request response for validation errors
func (b *backend) validationErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	b.errorResponse(w, r, http.StatusBadRequest, "VALIDATION_ERROR", message)
}

// conflictResponse sends a 409 Conflict response
func (b *backend) conflictResponse(w http.ResponseWriter, r *http.Request, message string) {
	b.errorResponse(w, r, http.StatusConflict, "CONFLICT", message)
}

// failedValidationResponse sends a 422 Unprocessable Entity response with validation errors
func (b *backend) failedValidationResponse(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	message := fmt.Sprintf("validation failed: %v", errs)
	b.errorResponse(w, r, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message)
}

func (b *backend) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s method not supported for this request", r.Method)
	b.errorResponse(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", message)
}
