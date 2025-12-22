package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

type TestBody struct {
	Str string  `json:"string"`
	Num int     `json:"int"`
	Err *string `json:"error,omitempty"`
}

type WriteSuccessTestCase struct {
	name           string
	writeFunc      func(w http.ResponseWriter, resp ResponseBody)
	respBody       ResponseBody
	expectedStatus int
}

func TestHttpHelper_WriteSuccess(t *testing.T) {
	msg := "Test Response"
	testData := TestBody{
		Str: msg,
		Num: 42,
		Err: nil,
	}

	testCases := []WriteSuccessTestCase{
		{
			name:      "Write Ok Response",
			writeFunc: WriteOkResponse,
			respBody: ResponseBody{
				Message: msg,
				Data:    testData,
				Error:   nil,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Write Created Response",
			writeFunc: WriteCreatedResponse,
			respBody: ResponseBody{
				Message: "created",
				Data:    testData,
				Error:   nil,
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Call the write function
			tc.writeFunc(w, tc.respBody)

			// Test status code
			assert.Equal(t, w.Code, tc.expectedStatus)

			// Test Content-Type header
			assert.Equal(t, w.Header().Get("Content-Type"), "application/json")

			// Test body contents
			var actualResponse ResponseBody
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NilError(t, err, "Failed to unmarshal response body")

			// Test message field
			assert.Equal(t, actualResponse.Message, tc.respBody.Message)

			// Test error field
			assert.DeepEqual(t, actualResponse.Error, tc.respBody.Error)
		})
	}
}

type WriteErrorTestCase struct {
	name           string
	writeFunc      func(w http.ResponseWriter, err error)
	errorObj       error
	expectedStatus int
}

func TestHttpHelper_WriteError(t *testing.T) {
	testCases := []WriteErrorTestCase{
		{
			name:           "Test Internal Error",
			writeFunc:      WriteInternalErrorResponse,
			errorObj:       fmt.Errorf("Test Internal"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Test Not Found Error",
			writeFunc:      WriteNotFoundResponse,
			errorObj:       fmt.Errorf("Test Not Found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Test Not Implemented Error",
			writeFunc:      WriteNotImplementedResponse,
			errorObj:       fmt.Errorf("Test Not Implemented"),
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name:           "Test Bad Request Error",
			writeFunc:      WriteBadRequestResponse,
			errorObj:       fmt.Errorf("Test Bad Request"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Test Unauthorized Error",
			writeFunc:      WriteUnauthorizedResponse,
			errorObj:       fmt.Errorf("Test Unauthorized"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Test Forbidden Error",
			writeFunc:      WriteForbiddenResponse,
			errorObj:       fmt.Errorf("Test Forbidden"),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Test Conflict Error",
			writeFunc:      WriteConflictResponse,
			errorObj:       fmt.Errorf("Test Conflict"),
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Call the write function
			tc.writeFunc(w, tc.errorObj)

			// Test status code
			assert.Equal(t, w.Code, tc.expectedStatus)

			// Test Content-Type header
			assert.Equal(t, w.Header().Get("Content-Type"), "application/json")

			// Test body contents
			var actualResponse ResponseBody
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NilError(t, err, "Failed to unmarshal response body")

			// Test message field
			assert.Equal(t, actualResponse.Message, "error")

			// Test error field
			assert.Equal(t, *actualResponse.Error, tc.errorObj.Error())
		})
	}
}
