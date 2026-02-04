package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockHandler is a simple handler that writes the value from context
type mockHandler struct {
	called  bool
	request *http.Request
}

// ServeHTTP implements the http.Handler interface for mockHandler
func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.called = true
	m.request = r
	w.WriteHeader(http.StatusOK)
}

type testCase struct {
	name           string
	pathValues     map[string]string
	valueNames     []string
	expectedStatus int
	expectNext     bool
}

// typeValidator is a function that validates the type of a value in the context
type typeValidator func(t *testing.T, valueName string, value any)

// runPathValueTests is a generic test runner for path value middleware
func runPathValueTests(
	t *testing.T,
	middlewareFunc func(...string) Middleware,
	tests []testCase,
	validator typeValidator,
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock handler
			mockNext := &mockHandler{}

			// Create the middleware
			middleware := middlewareFunc(tt.valueNames...)
			handler := middleware(mockNext)

			// Create a request with path values
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			for key, value := range tt.pathValues {
				req.SetPathValue(key, value)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Execute the handler
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, actual status %d", tt.expectedStatus, rr.Code)
			}

			// Check if next handler was called
			if mockNext.called != tt.expectNext {
				t.Errorf("expected next handler called: %v, got: %v", tt.expectNext, mockNext.called)
			}

			// If successful, verify the context contains the parsed values
			// and validator checks their types
			if tt.expectNext && tt.expectedStatus == http.StatusOK {
				if mockNext.request == nil {
					t.Fatal("expected request to be captured in mock handler")
				}
				for _, valueName := range tt.valueNames {
					value := mockNext.request.Context().Value(valueName)
					if value == nil {
						t.Errorf("expected value for %s to be in context", valueName)
						continue
					}
					validator(t, valueName, value)
				}
			}
		})
	}
}

func TestIntPathValues(t *testing.T) {
	tests := []testCase{
		// VALID CASES
		{
			name: "single valid integer",
			pathValues: map[string]string{
				"id": "123",
			},
			valueNames:     []string{"id"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "multiple valid integers",
			pathValues: map[string]string{
				"id":     "123",
				"userId": "456",
			},
			valueNames:     []string{"id", "userId"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "negative integer",
			pathValues: map[string]string{
				"id": "-42",
			},
			valueNames:     []string{"id"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "zero value",
			pathValues: map[string]string{
				"id": "0",
			},
			valueNames:     []string{"id"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "large integer",
			pathValues: map[string]string{
				"id": "9223372036854775807",
			},
			valueNames:     []string{"id"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		// INVALID CASES
		{
			name: "invalid integer",
			pathValues: map[string]string{
				"id": "abc",
			},
			valueNames:     []string{"id"},
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:       "missing path value",
			pathValues: map[string]string{
				// id is not provided
			},
			valueNames:     []string{"id"},
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name: "floating point as integer",
			pathValues: map[string]string{
				"id": "123.45",
			},
			valueNames:     []string{"id"},
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
	}

	validator := func(t *testing.T, valueName string, value any) {
		if _, ok := value.(int64); !ok {
			t.Errorf("expected value for %s to be int64, got %T", valueName, value)
		}
	}

	runPathValueTests(t, IntPathValues, tests, validator)
}

func TestFloatPathValues(t *testing.T) {
	tests := []testCase{
		// VALID CASES
		{
			name: "single valid float",
			pathValues: map[string]string{
				"price": "19.99",
			},
			valueNames:     []string{"price"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "multiple valid floats",
			pathValues: map[string]string{
				"price":  "19.99",
				"weight": "5.5",
			},
			valueNames:     []string{"price", "weight"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "integer as float",
			pathValues: map[string]string{
				"price": "42",
			},
			valueNames:     []string{"price"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "negative float",
			pathValues: map[string]string{
				"price": "-3.14",
			},
			valueNames:     []string{"price"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "scientific notation",
			pathValues: map[string]string{
				"price": "1.23e-4",
			},
			valueNames:     []string{"price"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "zero value",
			pathValues: map[string]string{
				"price": "0.0",
			},
			valueNames:     []string{"price"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		// INVALID CASES
		{
			name: "invalid float",
			pathValues: map[string]string{
				"price": "not-a-number",
			},
			valueNames:     []string{"price"},
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:       "missing path value",
			pathValues: map[string]string{
				// price is not provided
			},
			valueNames:     []string{"price"},
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
	}

	validator := func(t *testing.T, valueName string, value any) {
		if _, ok := value.(float64); !ok {
			t.Errorf("expected value for %s to be float64, got %T", valueName, value)
		}
	}

	runPathValueTests(t, FloatPathValues, tests, validator)
}

func TestStringPathValues(t *testing.T) {
	tests := []testCase{
		// VALID CASES
		{
			name: "single valid string",
			pathValues: map[string]string{
				"username": "john_doe",
			},
			valueNames:     []string{"username"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "multiple valid strings",
			pathValues: map[string]string{
				"username": "john_doe",
				"category": "food",
			},
			valueNames:     []string{"username", "category"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "string with special characters",
			pathValues: map[string]string{
				"slug": "hello-world-2024",
			},
			valueNames:     []string{"slug"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name: "numeric string",
			pathValues: map[string]string{
				"code": "12345",
			},
			valueNames:     []string{"code"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},

		{
			name: "unicode string",
			pathValues: map[string]string{
				"name": "café",
			},
			valueNames:     []string{"name"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		// INVALID CASES
		{
			name:       "missing path value",
			pathValues: map[string]string{
				// username is not provided
			},
			valueNames:     []string{"username"},
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name: "empty string value should fail",
			pathValues: map[string]string{
				"name": "",
			},
			valueNames:     []string{"name"},
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
	}

	validator := func(t *testing.T, valueName string, value any) {
		if _, ok := value.(string); !ok {
			t.Errorf("expected value for %s to be string, got %T", valueName, value)
		}
	}

	runPathValueTests(t, StringPathValues, tests, validator)
}
