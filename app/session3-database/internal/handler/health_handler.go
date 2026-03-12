package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
	db        *sql.DB
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		db:        db,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string        `json:"status"`
	Message   string        `json:"message"`
	Uptime    time.Duration `json:"uptime_seconds"`
	Timestamp time.Time     `json:"timestamp"`
}

// Check handles GET /health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {

	status := "ok"
	dbStatus := "connected"

	err := h.db.Ping()
	if err != nil {
		status = "degraded"
		dbStatus = "disconnected"
	}

	stats := h.db.Stats()

	response := map[string]interface{}{
		"status": status,
		"database": map[string]interface{}{
			"status":           dbStatus,
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
			"max_open":         stats.MaxOpenConnections,
		},
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")

	if status != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

/*
🎓 NOTES:

Refactored từ Buổi 1:
- Buổi 1: Health check logic trong main.go
- Buổi 2: Extracted to separate handler

Benefits:
- Consistent with other handlers
- Can add more health checks (database, etc.) in Buổi 3
- Reusable and testable
*/
