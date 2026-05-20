package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cainseing/drop-cli/internal/config"
)

func TestPostBlob(t *testing.T) {
	// Save original API URL
	originalURL := config.ApiURL
	defer func() { config.ApiURL = originalURL }()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/blob" {
			t.Errorf("Expected POST /blob, got %s %s", r.Method, r.URL.Path)
		}

		var req DropRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Check request fields
		if req.Blob != "test-blob" {
			t.Errorf("Expected blob 'test-blob', got %s", req.Blob)
		}
		if req.TTL != 300 { // 5 * 60
			t.Errorf("Expected TTL 300, got %d", req.TTL)
		}
		if req.Reads != 1 {
			t.Errorf("Expected reads 1, got %d", req.Reads)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "test-id"})
	}))
	defer server.Close()

	// Set API URL to test server
	config.ApiURL = server.URL

	// Call postBlob
	id, err := postBlob("test-blob", 5, 1, "sig", "sender", "provider")
	if err != nil {
		t.Errorf("postBlob failed: %v", err)
	}
	if id != "test-id" {
		t.Errorf("Expected id 'test-id', got %s", id)
	}
}

func TestPostBlob_Error(t *testing.T) {
	originalURL := config.ApiURL
	defer func() { config.ApiURL = originalURL }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Code:    400,
			Message: "Invalid request",
		})
	}))
	defer server.Close()

	config.ApiURL = server.URL

	_, err := postBlob("test-blob", 5, 1, "", "", "")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "request to API failed: Invalid request" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetBlob(t *testing.T) {
	originalURL := config.ApiURL
	defer func() { config.ApiURL = originalURL }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/blob/test-id" {
			t.Errorf("Expected GET /blob/test-id, got %s %s", r.Method, r.URL.Path)
		}

		response := GetDropResponse{
			Blob:           "test-blob",
			RemainingReads: 5,
			Signature:      "sig",
			Sender:         "sender",
			Provider:       "provider",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config.ApiURL = server.URL

	result, err := getBlob("test-id")
	if err != nil {
		t.Errorf("getBlob failed: %v", err)
	}
	if result.Blob != "test-blob" {
		t.Errorf("Expected blob 'test-blob', got %s", result.Blob)
	}
	if result.RemainingReads != 5 {
		t.Errorf("Expected remaining reads 5, got %d", result.RemainingReads)
	}
}

func TestGetBlob_NotFound(t *testing.T) {
	originalURL := config.ApiURL
	defer func() { config.ApiURL = originalURL }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config.ApiURL = server.URL

	_, err := getBlob("nonexistent")
	if err == nil {
		t.Error("Expected error for not found, got nil")
	}
	if err.Error() != "drop was not found" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestPurgeBlob(t *testing.T) {
	originalURL := config.ApiURL
	defer func() { config.ApiURL = originalURL }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" || r.URL.Path != "/blob/test-id" {
			t.Errorf("Expected DELETE /blob/test-id, got %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config.ApiURL = server.URL

	success, err := purgeBlob("test-id")
	if err != nil {
		t.Errorf("purgeBlob failed: %v", err)
	}
	if !success {
		t.Error("Expected success true, got false")
	}
}

func TestPurgeBlob_NotFound(t *testing.T) {
	originalURL := config.ApiURL
	defer func() { config.ApiURL = originalURL }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config.ApiURL = server.URL

	_, err := purgeBlob("nonexistent")
	if err == nil {
		t.Error("Expected error for not found, got nil")
	}
	if err.Error() != "drop was not found" {
		t.Errorf("Unexpected error: %v", err)
	}
}
