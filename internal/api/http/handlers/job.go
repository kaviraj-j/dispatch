package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kaviraj-j/dispatch/internal/queue"
)

type JobHandler struct {
	qm *queue.Manager
}

func NewJobHandler(qm *queue.Manager) *JobHandler {
	return &JobHandler{qm: qm}
}

// POST /jobs
func (h *JobHandler) HandleJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queueName := r.URL.Query().Get("queue")
	if queueName == "" {
		http.Error(w, "queue name is required", http.StatusBadRequest)
		return
	}

	payload := json.RawMessage{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// auto-create queue on enqueue
	q := h.qm.EnsureQueue(queueName)
	job := q.Enqueue(payload)

	json.NewEncoder(w).Encode(job)
}

// GET /jobs/consume?queue=name
func (h *JobHandler) HandleConsume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queueName := r.URL.Query().Get("queue")
	if queueName == "" {
		http.Error(w, "queue name is required", http.StatusBadRequest)
		return
	}

	q, ok := h.qm.GetQueue(queueName)
	if !ok {
		http.Error(w, "queue not found", http.StatusNotFound)
		return
	}

	job, ok := q.Dequeue()
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	json.NewEncoder(w).Encode(job)
}

// POST /jobs/ack?id=JOB_ID&queue=name
func (h *JobHandler) HandleAck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queueName := r.URL.Query().Get("queue")
	jobID := r.URL.Query().Get("id")

	if queueName == "" || jobID == "" {
		http.Error(w, "queue and job id are required", http.StatusBadRequest)
		return
	}

	q, ok := h.qm.GetQueue(queueName)
	if !ok {
		http.Error(w, "queue not found", http.StatusNotFound)
		return
	}

	if !q.Ack(jobID) {
		http.Error(w, "job not found or not in progress", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// POST /jobs/fail?id=JOB_ID&queue=name
func (h *JobHandler) HandleJobFail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queueName := r.URL.Query().Get("queue")
	jobID := r.URL.Query().Get("id")

	if queueName == "" || jobID == "" {
		http.Error(w, "queue and job id are required", http.StatusBadRequest)
		return
	}

	q, ok := h.qm.GetQueue(queueName)
	if !ok {
		http.Error(w, "queue not found", http.StatusNotFound)
		return
	}

	if !q.Ack(jobID) {
		http.Error(w, "job not found or not in progress", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /deadletter?queue=name
func (h *JobHandler) HandleDeadLetterQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queueName := r.URL.Query().Get("queue")

	if queueName == "" {
		http.Error(w, "queue and job id are required", http.StatusBadRequest)
		return
	}

	q, ok := h.qm.GetQueue(queueName)
	if !ok {
		http.Error(w, "queue not found", http.StatusNotFound)
		return
	}

	deadLetters := q.GetDeadLetterJobs()
	if len(deadLetters) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	json.NewEncoder(w).Encode(deadLetters)
}
