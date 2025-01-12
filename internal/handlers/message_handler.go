package handlers

import (
	"encoding/json"
	"message-service/internal/model/request"
	"net/http"
)

// CreateMessage godoc
// @Summary Create a new message
// @Description Create a new message with the given content
// @Tags messages
// @Accept json
// @Produce json
// @Param message body request.MessageRequest true "Message content"
// @Router /api/message [post]
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req request.MessageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.services.CreateMessage(req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"info": "Message saved!"})
}

// GetStats godoc
// @Summary Get message statistics
// @Description Get statistics of messages from the service
// @Tags messages
// @Produce json
// @Router /api/messages/stats [get]
func (h *Handler) GetStats(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	msgs, err := h.services.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(msgs)
}
