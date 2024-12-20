package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Client represents the video translator client
type Client struct {
	BaseURL    string        // Base URL of the server
	HTTPClient *http.Client  // HTTP client for requests
	PollDelay  time.Duration // Initial delay for polling
	MaxRetries int           // Max retries for polling
	logger     *logrus.Logger
}

// NewClient creates a new instance of the client
func NewClient(baseURL string) *Client {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		PollDelay:  1 * time.Second, // Default polling delay
		MaxRetries: 10,              // Default max retries
		logger:     logger,
	}
}

// StartJob starts a new translation job
func (c *Client) StartJob() (string, error) {
	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/start", c.BaseURL), "application/json", nil)
	if err != nil {
		c.logger.WithError(err).Error("Failed to start job")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		err := fmt.Errorf("failed to start job, status code: %d", resp.StatusCode)
		c.logger.WithError(err).Error("StartJob failed")
		return "", err
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.logger.WithError(err).Error("Failed to decode start job response")
		return "", err
	}

	jobID, ok := result["job_id"]
	if !ok {
		err := errors.New("invalid response: missing job_id")
		c.logger.WithError(err).Error("StartJob response missing job_id")
		return "", err
	}

	c.logger.WithField("job_id", jobID).Info("Job started successfully")
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
			c.logger.WithError(err).WithField("job_id", jobID).Error("Failed to get job status")
			return "", err
		}

		c.logger.WithFields(logrus.Fields{
			"job_id": jobID,
			"status": status,
			"retry":  retries,
		}).Info("Polled job status")

		if status == "completed" || status == "error" {
			return status, nil
		}

		// Exponential backoff
		delay *= 2
	}

	err := fmt.Errorf("job %s did not complete within the retry limit", jobID)
	c.logger.WithError(err).WithField("job_id", jobID).Error("PollStatus failed")
	return "", err
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
