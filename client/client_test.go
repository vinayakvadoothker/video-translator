package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStartJob(t *testing.T) {
	// Create a mock server to simulate /start endpoint
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/start" || r.Method != http.MethodPost {
			http.Error(w, "Invalid endpoint", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"job_id": "12345"}`))
	}))
	defer mockServer.Close()

	// Create a client pointing to the mock server
	client := NewClient(mockServer.URL)

	// Test StartJob
	jobID, err := client.StartJob()
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	if jobID != "12345" {
		t.Fatalf("Unexpected job ID: got %s, want %s", jobID, "12345")
	}
}

func TestPollStatus(t *testing.T) {
	// Create a mock server to simulate /status endpoint
	callCount := 0
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status" || r.Method != http.MethodGet {
			http.Error(w, "Invalid endpoint", http.StatusNotFound)
			return
		}

		callCount++
		w.Header().Set("Content-Type", "application/json")

		// Simulate different responses based on call count
		if callCount < 3 {
			w.Write([]byte(`{"status": "pending"}`))
		} else {
			w.Write([]byte(`{"status": "completed"}`))
		}
	}))
	defer mockServer.Close()

	// Create a client pointing to the mock server
	client := NewClient(mockServer.URL)
	client.PollDelay = 100 * time.Millisecond // Using shorter delay for faster testing
	client.MaxRetries = 5

	// Test PollStatus
	status, err := client.PollStatus("12345")
	if err != nil {
		t.Fatalf("Failed to poll status: %v", err)
	}

	if status != "completed" {
		t.Fatalf("Unexpected status: got %s, want %s", status, "completed")
	}
}
