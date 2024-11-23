package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client represents the video translator client
type Client struct {
	BaseURL    string        // Base URL of the server
	HTTPClient *http.Client  // HTTP client for requests
	PollDelay  time.Duration // Initial delay for polling
	MaxRetries int           // Max retries for polling
}

// NewClient creates a new instance of the client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		PollDelay:  1 * time.Second, // Default polling delay
		MaxRetries: 10,              // Default max retries
	}
}

// StartJob starts a new translation job
func (c *Client) StartJob() (string, error) {
	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/start", c.BaseURL), "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to start job, status code: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	jobID, ok := result["job_id"]
	if !ok {
		return "", errors.New("invalid response: missing job_id")
	}

	return jobID, nil
}

// PollStatus polls the status of a job until it is completed or errors out
func (c *Client) PollStatus(jobID string) (string, error) {
	var retries int
	delay := c.PollDelay

	for retries = 0; retries < c.MaxRetries; retries++ {
		time.Sleep(delay)

		status, err := c.getStatus(jobID)
		if err != nil {
			return "", err
		}

		if status == "completed" || status == "error" {
			return status, nil
		}

		// Exponential backoff
		delay *= 2
	}

	return "", fmt.Errorf("job %s did not complete within the retry limit", jobID)
}

// getStatus gets the current status of a job
func (c *Client) getStatus(jobID string) (string, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/status?job_id=%s", c.BaseURL, jobID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch status, status code: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	status, ok := result["status"]
	if !ok {
		return "", errors.New("invalid response: missing status")
	}

	return status, nil
}
