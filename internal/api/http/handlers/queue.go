package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kaviraj-j/dispatch/internal/queue"
)

type QueueHandler struct {
	qm *queue.Manager
}

func NewQueueHandler(qm *queue.Manager) *QueueHandler {
	return &QueueHandler{qm: qm}
}

// POST /queues
// GET  /queues
func (h *QueueHandler) HandleQueues(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createQueue(w, r)
	case http.MethodGet:
		h.listQueues(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

type createQueueRequest struct {
	Name string `json:"name"`
}

func (h *QueueHandler) createQueue(w http.ResponseWriter, r *http.Request) {
	var req createQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err := h.qm.CreateQueue(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *QueueHandler) listQueues(w http.ResponseWriter, r *http.Request) {
	queues := h.qm.ListQueues()
	json.NewEncoder(w).Encode(queues)
}
