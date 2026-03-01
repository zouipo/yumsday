package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name           string
		appErr         *AppError
		expectedOutput string
	}{
		{
			name: "with underlying error",
			appErr: &AppError{
				Message:    "Failed to fetch user",
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("database connection failed"),
			},
			expectedOutput: "Failed to fetch user: database connection failed",
		},
		{
			name: "without underlying error",
			appErr: &AppError{
				Message:    "User not found",
				StatusCode: http.StatusNotFound,
				Err:        nil,
			},
			expectedOutput: "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.appErr.Error()
			if actual != tt.expectedOutput {
				t.Errorf("Error() = %q, expected %q", actual, tt.expectedOutput)
			}
		})
	}
}

func TestErrorConstructors(t *testing.T) {
	underlyingErr := errors.New("test error")

	tests := []struct {
		name            string
		createError     func() error
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "NewEntityNotFoundError",
			createError:     func() error { return NewEntityNotFoundError("User", "123", underlyingErr) },
			expectedStatus:  http.StatusNotFound,
			expectedMessage: "User 123 not found",
		},
		{
			name:            "NewValidationError",
			createError:     func() error { return NewValidationError("username", "must be at least 3 characters", underlyingErr) },
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Validation error on field 'username': must be at least 3 characters",
		},
		{
			name:            "NewUnauthorizedError",
			createError:     func() error { return NewUnauthorizedError(underlyingErr) },
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "Unauthorized: authentication required",
		},
		{
			name:            "NewForbiddenError",
			createError:     func() error { return NewForbiddenError(underlyingErr) },
			expectedStatus:  http.StatusForbidden,
			expectedMessage: "Forbidden: user not allowed to perform this action",
		},
		{
			name:            "NewConflictError",
			createError:     func() error { return NewConflictError("User", "username already exists", underlyingErr) },
			expectedStatus:  http.StatusConflict,
			expectedMessage: "Conflict with User: username already exists",
		},
		{
			name:            "NewInternalServerError",
			createError:     func() error { return NewInternalServerError("Failed to process request", underlyingErr) },
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed to process request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createError()

			appErr, ok := err.(*AppError)
			if !ok {
				t.Fatal("Expected error to be of type *AppError")
			}

			if appErr.StatusCode != tt.expectedStatus {
				t.Errorf("StatusCode = %d, expected %d", appErr.StatusCode, tt.expectedStatus)
			}

			if appErr.Message != tt.expectedMessage {
				t.Errorf("Message = %q, expected %q", appErr.Message, tt.expectedMessage)
			}

			if appErr.Err != underlyingErr {
				t.Errorf("Underlying error not preserved")
			}

			expectedError := tt.expectedMessage + ": " + underlyingErr.Error()
			if err.Error() != expectedError {
				t.Errorf("Error() = %q, expected %q", err.Error(), expectedError)
			}
		})
	}
}

func TestNewInternalServerError_WithNilError(t *testing.T) {
	err := NewInternalServerError("Something went wrong", nil)

	appErr, ok := err.(*AppError)
	if !ok {
		t.Fatal("Expected error to be of type *AppError")
	}

	if appErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, expected %d", appErr.StatusCode, http.StatusInternalServerError)
	}

	expectedMessage := "Something went wrong"
	if appErr.Message != expectedMessage {
		t.Errorf("Message = %q, expected %q", appErr.Message, expectedMessage)
	}

	if appErr.Err != nil {
		t.Errorf("Expected nil underlying error")
	}

	if err.Error() != expectedMessage {
		t.Errorf("Error() = %q, expected %q", err.Error(), expectedMessage)
	}
}
