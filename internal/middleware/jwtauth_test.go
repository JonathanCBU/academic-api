package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

// Test constants
const (
	AuthHeader   string = "Authorization"
	AuthPrefix   string = "Bearer"
	InvalidToken string = "1234"
	ValidToken   string = "12345"
	TestUrl      string = "/api/test/"
)

type ExtractAuthTestCase struct {
	name             string
	header           string
	expectedResult   string
	expectedErrorMsg string
}

func TestJwtAuth_extractAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testsCases := []ExtractAuthTestCase{
		{
			name:             "Successful header parse",
			header:           fmt.Sprintf("%s %s", AuthPrefix, ValidToken),
			expectedResult:   ValidToken,
			expectedErrorMsg: "",
		},
		{
			name:             "Auth header not present",
			header:           "",
			expectedResult:   "",
			expectedErrorMsg: "Auth header not present.",
		},
		{
			name:             "Auth header malformed",
			header:           ValidToken,
			expectedResult:   "",
			expectedErrorMsg: "Malformed auth header.",
		},
	}

	mw := NewJwtMiddleware(AuthHeader, AuthPrefix)

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up router
			router := mux.NewRouter()
			router.HandleFunc(TestUrl, func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)

			// Set up request
			req, _ := http.NewRequest(http.MethodGet, TestUrl, nil)
			if tc.header != "" {
				req.Header.Add(AuthHeader, tc.header)
			}

			// Write asserts in middleware handler
			router.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					resp, err := mw.extractAuth(req.Header)
					if tc.expectedErrorMsg != "" {
						assert.Error(t, err, tc.expectedErrorMsg)
					} else {
						assert.NilError(t, err)
					}
					assert.Equal(t, resp, tc.expectedResult)
					next.ServeHTTP(w, req)
				})
			})

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
		})
	}
}

type CheckAuthTestCase struct {
	name             string
	token            string
	expectedResult   bool
	expectedErrorMsg string
}

func TestJwtAuth_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testsCases := []CheckAuthTestCase{
		{
			name:             "Valid token",
			token:            ValidToken,
			expectedResult:   true,
			expectedErrorMsg: "",
		},
		{
			name:             "Invalid token",
			token:            InvalidToken,
			expectedResult:   false,
			expectedErrorMsg: "Token min length 5",
		},
	}

	mw := NewJwtMiddleware(AuthHeader, AuthPrefix)

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := mw.CheckAuth(tc.token)
			if tc.expectedErrorMsg != "" {
				assert.Error(t, err, tc.expectedErrorMsg)
			} else {
				assert.NilError(t, err)
			}
			assert.Equal(t, resp, tc.expectedResult)
		})
	}
}
