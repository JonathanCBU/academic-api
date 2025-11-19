package handler

import (
	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	Message string  `json:"message"`
	Data    any     `json:"data"`
	Error   *string `json:"error,omitempty"`
}

// WriteHttpResponse writes a JSON response with the given status code
func WriteHttpResponse(w http.ResponseWriter, apiResponse ResponseBody, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	responsePayload, err := json.Marshal(apiResponse)
	if err != nil {
		// Fallback error response
		errorPayload, _ := json.Marshal(map[string]string{"error": "Failed to marshal JSON response"})
		_, _ = w.Write(errorPayload)
		return
	}
	_, _ = w.Write(responsePayload)
}

// WriteOkResponse writes a 200 OK response
func WriteOkResponse(w http.ResponseWriter, respBody ResponseBody) {
	WriteHttpResponse(w, respBody, http.StatusOK)
}

// WriteCreatedResponse writes a 201 Created response
func WriteCreatedResponse(w http.ResponseWriter, respBody ResponseBody) {
	WriteHttpResponse(w, respBody, http.StatusCreated)
}

// WriteErrorResponse writes an error response with the given status code
func WriteErrorResponse(w http.ResponseWriter, err error, httpStatusCode int) {
	errorMessage := err.Error()
	respBody := ResponseBody{
		Message: "error",
		Data:    nil,
		Error:   &errorMessage,
	}
	WriteHttpResponse(w, respBody, httpStatusCode)
}

// WriteInternalErrorResponse writes a 500 Internal Server Error response
func WriteInternalErrorResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, err, http.StatusInternalServerError)
}

// WriteNotFoundResponse writes a 404 Not Found response
func WriteNotFoundResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, err, http.StatusNotFound)
}

// WriteNotImplementedResponse writes a 501 Not Implemented response
func WriteNotImplementedResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, err, http.StatusNotImplemented)
}

// WriteBadRequestResponse writes a 400 Bad Request response
func WriteBadRequestResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, err, http.StatusBadRequest)
}

// WriteUnauthorizedResponse writes a 401 Unauthorized response
func WriteUnauthorizedResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, err, http.StatusUnauthorized)
}

// WriteForbiddenResponse writes a 403 Forbidden response
func WriteForbiddenResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, err, http.StatusForbidden)
}

// WriteConflictResponse writes a 409 Conflict response
func WriteConflictResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, err, http.StatusConflict)
}
