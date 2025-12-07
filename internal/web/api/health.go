package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	redisClient *redis.Client
	db          *pgxpool.Pool
	logger      *slog.Logger
}

type DetailedHealthResponse struct {
	Status    string            `json:"status"`
	Checks    map[string]string `json:"checks"`
	Timestamp time.Time         `json:"timestamp"`
}

func NewHealthHandler(redisClient *redis.Client, db *pgxpool.Pool, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		redisClient: redisClient,
		db:          db,
		logger:      logger,
	}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	checks := make(map[string]string)
	allHealthy := true

	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		checks["redis"] = "down"
		allHealthy = false
	} else {
		checks["redis"] = "up"
	}

	if h.db != nil {
		if err := h.db.Ping(ctx); err != nil {
			checks["database"] = "down"
			allHealthy = false
		} else {
			checks["database"] = "up"
		}
	}

	response := DetailedHealthResponse{
		Status:    "healthy",
		Checks:    checks,
		Timestamp: time.Now(),
	}

	statusCode := http.StatusOK
	if !allHealthy {
		response.Status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	h.HealthCheck(w, r)
}
