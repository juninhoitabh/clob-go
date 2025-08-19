package httpClient_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpClient "github.com/juninhoitabh/clob-go/internal/infra/http-client"
)

func TestNewDefaultHttpClient(t *testing.T) {
	timeout := 5 * time.Second
	client := httpClient.NewDefaultHttpClient(timeout)
	assert.NotNil(t, client, "O cliente HTTP não deve ser nulo")
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		assert.Equal(t, "TestValue", r.Header.Get("TestHeader"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()

	client := httpClient.NewDefaultHttpClient(5 * time.Second)
	headers := map[string]string{"TestHeader": "TestValue"}

	resp, err := client.Get(t.Context(), server.URL, headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(resp.Body), "success")
	assert.Equal(t, "application/json", resp.Headers.Get("Content-Type"))
}

func TestPost(t *testing.T) {
	type testPayload struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	payload := testPayload{Name: "test", Value: 42}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "TestValue", r.Header.Get("TestHeader"))

		var receivedPayload testPayload
		err := json.NewDecoder(r.Body).Decode(&receivedPayload)
		assert.NoError(t, err)
		assert.Equal(t, payload.Name, receivedPayload.Name)
		assert.Equal(t, payload.Value, receivedPayload.Value)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"123"}`))
	}))
	defer server.Close()

	client := httpClient.NewDefaultHttpClient(5 * time.Second)
	headers := map[string]string{"TestHeader": "TestValue"}

	resp, err := client.Post(t.Context(), server.URL, payload, headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Contains(t, string(resp.Body), "123")
	assert.Equal(t, "application/json", resp.Headers.Get("Content-Type"))
}

func TestPut(t *testing.T) {
	type testPayload struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	payload := testPayload{Name: "updated", Value: 100}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var receivedPayload testPayload
		err := json.NewDecoder(r.Body).Decode(&receivedPayload)
		assert.NoError(t, err)
		assert.Equal(t, payload.Name, receivedPayload.Name)
		assert.Equal(t, payload.Value, receivedPayload.Value)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"updated"}`))
	}))
	defer server.Close()

	client := httpClient.NewDefaultHttpClient(5 * time.Second)

	resp, err := client.Put(t.Context(), server.URL, payload, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(resp.Body), "updated")
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := httpClient.NewDefaultHttpClient(5 * time.Second)

	resp, err := client.Delete(t.Context(), server.URL, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server error"}`))
	}))
	defer server.Close()

	client := httpClient.NewDefaultHttpClient(5 * time.Second)

	resp, err := client.Get(t.Context(), server.URL, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, string(resp.Body), "server error")
}

func TestInvalidURL(t *testing.T) {
	client := httpClient.NewDefaultHttpClient(5 * time.Second)

	_, err := client.Get(t.Context(), "http://invalid-url-that-does-not-exist.xyz", nil)

	assert.Error(t, err)
}

func TestTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := httpClient.NewDefaultHttpClient(100 * time.Millisecond)

	_, err := client.Get(t.Context(), server.URL, nil)

	assert.Error(t, err)
}
