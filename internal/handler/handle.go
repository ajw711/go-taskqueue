package handler

import (
	"encoding/json"
	"go-taskqueue/internal/job"
	"go-taskqueue/internal/queue"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	queue *queue.Queue
}

type submitRequest struct {
	Type        string `json:"type"`
	Payload     string `json:"payload"`
	Priority    string `json:"priority"`
	MaxAttempts int    `json:"max_attempts"`
}

func NewHandler(q *queue.Queue) *Handler {
	return &Handler{queue: q}
}

func (h *Handler) SubmitJob(w http.ResponseWriter, r *http.Request) {
	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청", http.StatusBadRequest)
		return
	}

	priority := job.Priority(req.Priority)
	
	j, err := job.NewJob(uuid.New().String(), req.Type, req.Payload, priority, req.MaxAttempts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.queue.Submit(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
}


func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	j, flag := h.queue.GetJob(id)
	if flag == false {
		http.Error(w, "해당 job을 찾을 수 없습니다.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
}

func (h *Handler) DeadLetters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.queue.ListDeadLetters())
}