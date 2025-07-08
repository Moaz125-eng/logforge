package server

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	NodeID    string    `json:"node_id"`
	Timestamp time.Time `json:"timestamp"`
}

func HealthHandler(nodeID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HealthResponse{
			Status:    "ok",
			Service:   "logforge",
			NodeID:    nodeID,
			Timestamp: time.Now().UTC(),
		})
	}
}
