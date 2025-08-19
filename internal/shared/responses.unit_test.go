package shared_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/juninhoitabh/clob-go/internal/shared"
)

func TestMul(t *testing.T) {
	testCases := []struct {
		name     string
		a        int64
		b        int64
		expected int64
	}{
		{
			name:     "positive numbers",
			a:        2,
			b:        3,
			expected: 6,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -3,
			expected: 6,
		},
		{
			name:     "mixed sign",
			a:        -2,
			b:        3,
			expected: -6,
		},
		{
			name:     "zero",
			a:        0,
			b:        5,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := shared.Mul(tc.a, tc.b)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestWriteJSON(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	testData := testStruct{
		Name:  "test",
		Value: 42,
	}

	w := httptest.NewRecorder()

	shared.WriteJSON(w, http.StatusOK, testData)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response testStruct
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, testData.Name, response.Name)
	assert.Equal(t, testData.Value, response.Value)
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()

	err := errors.New("test error")
	shared.WriteError(w, err, http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.ErrorResponse
	decodeErr := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, decodeErr)
	assert.Equal(t, "test error", response.Message)
	assert.Equal(t, http.StatusBadRequest, response.Status)
}

func TestHandleError(t *testing.T) {
	testCases := []struct {
		err            error
		name           string
		expectedStatus int
	}{
		{
			name:           "ErrNotFound",
			err:            shared.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrAlreadyExists",
			err:            shared.ErrAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "ErrInvalidParam",
			err:            shared.ErrInvalidParam,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "default error",
			err:            errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			shared.HandleError(w, tc.err)
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response shared.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tc.err.Error(), response.Message)
			assert.Equal(t, tc.expectedStatus, response.Status)
		})
	}
}

func TestBadRequestError(t *testing.T) {
	testCases := []struct {
		name     string
		message  string
		details  []string
		expected shared.ErrorResponse
	}{
		{
			name:    "with message only",
			message: "Invalid request",
			details: nil,
			expected: shared.ErrorResponse{
				Message: "Invalid request",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name:    "with message and details",
			message: "Validation failed",
			details: []string{"Field 'name' is required", "Field 'email' must be valid"},
			expected: shared.ErrorResponse{
				Message: "Validation failed",
				Details: []string{"Field 'name' is required", "Field 'email' must be valid"},
				Status:  http.StatusBadRequest,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			shared.BadRequestError(w, tc.message, tc.details...)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response shared.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tc.expected.Message, response.Message)
			assert.Equal(t, tc.expected.Status, response.Status)
			assert.Equal(t, tc.expected.Details, response.Details)
		})
	}
}
