package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vinayakvadoothker/video-translator/client"
	"github.com/vinayakvadoothker/video-translator/server"
)

func TestIntegration(t *testing.T) {
	// Start the server in-memory
	srv := httptest.NewServer(http.HandlerFunc(server.Router)) // server.Router must be the main router of your server
	defer srv.Close()

	// Create a client instance pointing to the test server
	client := client.NewClient(srv.URL)

	// Start a new job
	jobID, err := client.StartJob()
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}
	t.Logf("Started job with ID: %s", jobID)

	// Poll for status
	status, err := client.PollStatus(jobID)
	if err != nil {
		t.Fatalf("Failed to poll job status: %v", err)
	}
	t.Logf("Job %s completed with status: %s", jobID, status)

	// Assert the final status
	if status != "completed" && status != "error" {
		t.Fatalf("Unexpected status: %s", status)
	}
}
