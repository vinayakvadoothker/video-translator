package client

import (
	"testing"
)

func TestClient(t *testing.T) {
	baseURL := "http://localhost:8080" // The server is running locally

	client := NewClient(baseURL)

	// Test starting a job
	jobID, err := client.StartJob()
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}
	t.Logf("Started job with ID: %s", jobID)

	// Test polling for status
	status, err := client.PollStatus(jobID)
	if err != nil {
		t.Fatalf("Failed to poll status: %v", err)
	}
	t.Logf("Job %s completed with status: %s", jobID, status)

	// Makes sure status is valid
	if status != "completed" && status != "error" {
		t.Fatalf("Unexpected status: %s", status)
	}
}
