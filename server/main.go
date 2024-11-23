package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Job struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	StartedAt time.Time
	Duration  time.Duration
}

var (
	jobs       = make(map[string]*Job) // Store jobs in memory
	jobsMutex  = sync.Mutex{}          // Mutex for concurrent access
	jobTimeout time.Duration           // Configurable timeout for completion
)

func main() {
	// Read delay from environment variable (default: 10s)
	delay, err := strconv.Atoi(os.Getenv("JOB_TIMEOUT"))
	if err != nil {
		delay = 10 // Default delay is 10 seconds
	}
	jobTimeout = time.Duration(delay) * time.Second

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/start", handleStartJob)

	port := ":8080"
	log.Printf("Server is running on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Start a new job
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

// Check the status of a job
func handleStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job_id")
	if jobID == "" {
		http.Error(w, "Missing job_id query parameter", http.StatusBadRequest)
		return
	}

	jobsMutex.Lock()
	job, exists := jobs[jobID]
	if !exists {
		jobsMutex.Unlock()
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID, "status": job.Status})
}
