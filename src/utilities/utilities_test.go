package utilities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUUIDFromRequest(t *testing.T) {
	t.Run("Success Case", func(t *testing.T) {
		// Arrange
		expectedUUID := uuid.New()
		r := httptest.NewRequest(http.MethodGet, "/example", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", expectedUUID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		// Act
		parsedUUID, err := ParseUUIDFromRequest(r, "id")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUUID, parsedUUID)
	})

	t.Run("Invalid UUID Format", func(t *testing.T) {
		// Arrange
		r := httptest.NewRequest(http.MethodGet, "/example", nil)
		rctx := chi.NewRouteContext()

		rctx.URLParams.Add("id", "invalid-uuid")

		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		// Act
		parsedUUID, err := ParseUUIDFromRequest(r, "id")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, parsedUUID)
	})
}

func TestWriteJSONResponse(t *testing.T) {
	t.Run("Success Case", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		data := map[string]string{"message": "success"}

		// Act
		_ = WriteJSONResponse(w, data)

		var responseData map[string]string
		err := json.NewDecoder(w.Body).Decode(&responseData)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		assert.NoError(t, err)
		assert.Equal(t, data, responseData)
	})

	t.Run("Invalid Data for JSON Encoding", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		data := make(chan int)

		// Act
		err := WriteJSONResponse(w, data)
		if err != nil {
			t.Logf("invalid data for JSON encoding: %v", err)
		} else {
			t.Errorf("expected an error when encoding invalid data, but got none")
		}

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestReadRequestBody(t *testing.T) {
	t.Run("Success Case", func(t *testing.T) {
		// Arrange
		requestBody := []byte("example request body")
		r := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(requestBody))

		// Act
		bodyBytes, err := ReadRequestBody(r)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, requestBody, bodyBytes)
	})

	t.Run("Error Reading Request Body", func(t *testing.T) {
		// Arrange
		r := httptest.NewRequest(http.MethodPost, "/example", io.NopCloser(bytes.NewReader(nil)))
		r.Body = io.NopCloser(&readError{})

		// Act
		bodyBytes, err := ReadRequestBody(r)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, bodyBytes)
		assert.Contains(t, err.Error(), "failed to read request body")
	})
}

type readError struct{}

func (r *readError) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("read error")
}
