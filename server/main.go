package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Job struct {
	ID        string        `json:"id"`
	Status    string        `json:"status"`
	StartedAt time.Time     `json:"started_at"`
	Duration  time.Duration `json:"duration"`
}

var (
	logger     = logrus.New()          // Global logger instance
	jobs       = make(map[string]*Job) // Store jobs in memory
	jobsMutex  = sync.Mutex{}          // Mutex for thread-safe access
	jobTimeout time.Duration           // Configurable timeout for completion
)

// Router sets up and returns the HTTP routes for the server
func Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", handleStatus)
	mux.HandleFunc("/start", handleStartJob)
	return mux
}

// Start initializes the server on the specified port
func Start(port string) error {
	// Configure logger
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Read job timeout from environment variable (default: 10s)
	delay, err := strconv.Atoi(os.Getenv("JOB_TIMEOUT"))
	if err != nil {
		logger.WithError(err).Warn("Invalid JOB_TIMEOUT, using default 10 seconds")
		delay = 10
	}
	jobTimeout = time.Duration(delay) * time.Second

	logger.WithField("port", port).Info("Starting server")
	return http.ListenAndServe(":"+port, Router())
}

// handleStartJob starts a new job and returns the job ID
func handleStartJob(w http.ResponseWriter, r *http.Request) {
	jobID := strconv.Itoa(rand.Int())
	job := &Job{
		ID:        jobID,
		Status:    "pending",
		StartedAt: time.Now(),
		Duration:  jobTimeout,
	}

	jobsMutex.Lock()
	jobs[jobID] = job
	jobsMutex.Unlock()

	logger.WithField("job_id", jobID).Info("Job started")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

// handleStatus checks the status of a job and returns it
func handleStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job_id")
	if jobID == "" {
		logger.Warn("Missing job_id query parameter")
		http.Error(w, "Missing job_id query parameter", http.StatusBadRequest)
		return
	}

	jobsMutex.Lock()
	job, exists := jobs[jobID]
	if !exists {
		jobsMutex.Unlock()
		logger.WithField("job_id", jobID).Warn("Job not found")
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Simulate job lifecycle
	if time.Since(job.StartedAt) > job.Duration {
		job.Status = "completed"
	}
	if rand.Float32() < 0.05 { // Simulate a 5% chance of error
		job.Status = "error"
	}

	jobsMutex.Unlock()

	logger.WithFields(logrus.Fields{
		"job_id": jobID,
		"status": job.Status,
	}).Info("Job status checked")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID, "status": job.Status})
}
